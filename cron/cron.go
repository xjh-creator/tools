package cron

import (
	"context"
	"sort"
	"sync"
	"time"
)

type Cron struct {
	entries      []*Entry          //任务列表
	stop         chan struct{}     //停止
	add          chan *Entry       //添加
	edit         chan *Entry       //修改
	remove       chan EntryID      //删除
	snapshot     chan chan []Entry //快照
	logger       Logger            //日志
	running      bool              //是否运行
	runningMutex sync.Mutex        //运行锁
	location     *time.Location    //时区
	jobWaiter    sync.WaitGroup    //每个任务的运行状态
}

func New() *Cron {
	c := &Cron{
		entries:      nil,
		add:          make(chan *Entry),
		edit:         make(chan *Entry),
		stop:         make(chan struct{}),
		snapshot:     make(chan chan []Entry),
		remove:       make(chan EntryID),
		running:      false,
		runningMutex: sync.Mutex{},
		logger:       DefaultLogger,
		location:     time.Local,
	}
	currentTime, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		c.location = currentTime
	}
	return c
}

//添加任务
func (c *Cron) AddSchedule(schedule Schedule) EntryID {
	c.runningMutex.Lock()
	defer c.runningMutex.Unlock()
	entry := &Entry{
		ID:       EntryID(schedule.GetID()),
		Schedule: schedule,
	}
	if c.running {
		if c.IsRuning(entry.ID) {
			c.edit <- entry
		} else {
			c.add <- entry
		}
	} else {
		if c.IsRuning(entry.ID) {
			c.removeEntry(entry.ID)
		}
		c.entries = append(c.entries, entry)
	}
	return entry.ID
}

//IsRuning 任务是否已存在
func (c *Cron) IsRuning(ID EntryID) bool {
	for _, e := range c.entries {
		if e.ID == ID {
			return true
		}
	}
	return false
}

//Entries 返回任务列表
func (c *Cron) Entries() []Entry {
	c.runningMutex.Lock()
	defer c.runningMutex.Unlock()
	if c.running {
		replyChan := make(chan []Entry, 1)
		c.snapshot <- replyChan
		return <-replyChan
	}
	return c.entrySnapshot()
}

//Location 获取当前时区
func (c *Cron) Location() *time.Location {
	return c.location
}

//GetEntry 获取1个任务
func (c *Cron) GetEntry(id EntryID) Entry {
	for _, entry := range c.Entries() {
		if id == entry.ID {
			return entry
		}
	}
	return Entry{}
}

//RemoveEntry 删除任务
func (c *Cron) RemoveEntry(id EntryID) {
	c.runningMutex.Lock()
	defer c.runningMutex.Unlock()
	if c.running {
		c.remove <- id
	} else {
		c.removeEntry(id)
	}
}

//Start 协程运行定时任务
func (c *Cron) Start() {
	c.runningMutex.Lock()
	defer c.runningMutex.Unlock()
	if c.running {
		return
	}
	c.running = true
	go c.run()
}

//BlockRun 阻塞运行定时任务
func (c *Cron) BlockRun() {
	c.runningMutex.Lock()
	if c.running {
		c.runningMutex.Unlock()
		return
	}
	c.running = true
	c.runningMutex.Unlock()
	c.run()
}

//run 运行任务调度
func (c *Cron) run() {
	c.logger.Info("开始运行")

	//计算所有任务下次调度时间
	now := c.now()
	for _, entry := range c.entries {
		entry.Next = entry.Schedule.Next(now)
		c.logger.Info("调度", "现在", now, "任务ID", entry.ID, "下一次", entry.Next)
	}

	for {
		//按下一次执行顺序排序
		sort.Sort(byTime(c.entries))

		var timer *time.Timer
		if len(c.entries) == 0 || c.entries[0].Next.IsZero() {
			//没有可执行任务,休眠
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			timer = time.NewTimer(c.entries[0].Next.Sub(now))
		}

		for {
			select {
			case now = <-timer.C:
				now = now.In(c.location)
				c.logger.Info("唤醒", "时间", now)

				//循环运行所有的可执行任务
				for _, e := range c.entries {
					if e.Next.After(now) || e.Next.IsZero() {
						break
					}
					c.startJob(e, now)
					e.Prev = e.Next
					e.Next = e.Schedule.Next(now)
					c.logger.Info("运行", "现在", now, "任务ID", e.ID, "下一次", e.Next)
				}

			case newEntry := <-c.add:
				timer.Stop()
				now = c.now()
				newEntry.Next = newEntry.Schedule.Next(now)
				c.entries = append(c.entries, newEntry)
				c.logger.Info("添加", "现在", now, "任务ID", newEntry.ID, "下一次", newEntry.Next)
			case newEntry := <-c.edit:
				timer.Stop()
				now = c.now()
				c.removeEntry(newEntry.ID)
				newEntry.Next = newEntry.Schedule.Next(now)
				c.entries = append(c.entries, newEntry)
				c.logger.Info("修改", "现在", now, "任务ID", newEntry.ID, "下一次", newEntry.Next)

			case id := <-c.remove:
				timer.Stop()
				now = c.now()
				c.removeEntry(id)
				c.logger.Info("移除", "任务ID", id)

			case replyChan := <-c.snapshot:
				replyChan <- c.entrySnapshot()
				continue

			case <-c.stop:
				timer.Stop()
				c.logger.Info("停止")
				return
			}

			break
		}
	}
}

//协程运行任务
func (c *Cron) startJob(j *Entry, now time.Time) {
	c.jobWaiter.Add(1)
	go func() {
		defer c.jobWaiter.Done()
		j.Schedule.Work(now)
	}()
}

//返回时区的当前时间
func (c *Cron) now() time.Time {
	return time.Now().In(c.location)
}

//停止运行
func (c *Cron) Stop() context.Context {
	c.runningMutex.Lock()
	defer c.runningMutex.Unlock()
	if c.running {
		c.stop <- struct{}{}
		c.running = false
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c.jobWaiter.Wait()
		cancel()
	}()
	return ctx
}

// entrySnapshot returns a copy of the current cron entry list.
func (c *Cron) entrySnapshot() []Entry {
	var entries = make([]Entry, len(c.entries))
	for i, e := range c.entries {
		entries[i] = *e
	}
	return entries
}

//删除一个任务
func (c *Cron) removeEntry(id EntryID) {
	var entries []*Entry
	for _, e := range c.entries {
		if e.ID != id {
			entries = append(entries, e)
		}
	}
	c.entries = entries
}

//任务
type Schedule interface {
	//Next 获取下次任务
	Next(time.Time) time.Time
	//Work 执行任务
	Work(time.Time)
	//GetID 获取任务UUID
	GetID() string
	//GetSort 排序,用于时间一致时,越大越先执行
	GetSort() int64
}

type EntryID string

//Entry 任务
type Entry struct {
	ID EntryID
	//任务
	Schedule Schedule
	//下次运行时间
	Next time.Time
	//最后次运行时间
	Prev time.Time
	//开始时间
	BeginTime time.Time
	//结束时间
	StopTime time.Time
}

//是否非法的任务
func (e Entry) Valid() bool {
	return e.ID != ""
}

//byTime 任务排序
type byTime []*Entry

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	if s[i].Next.IsZero() {
		return false
	}
	if s[j].Next.IsZero() {
		return true
	}
	if s[i].Next.Equal(s[j].Next) {
		return s[i].Schedule.GetSort() > s[i].Schedule.GetSort()
	}
	return s[i].Next.Before(s[j].Next)
}
