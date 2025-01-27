package cointop

import (
	"strings"
	"time"
)

func (ct *Cointop) refresh() error {
	go func() {
		<-ct.limiter
		ct.forceRefresh <- true
	}()
	return nil
}

func (ct *Cointop) refreshAll() error {
	ct.refreshMux.Lock()
	defer ct.refreshMux.Unlock()
	ct.setRefreshStatus()
	ct.cache.Delete("allCoinsSlugMap")
	ct.cache.Delete("market")
	go func() {
		ct.updateCoins()
		ct.updateTable()
	}()
	go ct.UpdateChart()
	return nil
}

func (ct *Cointop) setRefreshStatus() {
	go func() {
		ct.loadingTicks("refreshing", 900)
		ct.rowChanged()
	}()
}

func (ct *Cointop) loadingTicks(s string, t int) {
	interval := 150
	k := 0
	for i := 0; i < (t / interval); i++ {
		ct.updateStatusbar(s + strings.Repeat(".", k))
		time.Sleep(time.Duration(i*interval) * time.Millisecond)
		k = k + 1
		if k > 3 {
			k = 0
		}
	}
}

func (ct *Cointop) intervalFetchData() {
	go func() {
		for {
			select {
			case <-ct.forceRefresh:
				ct.refreshAll()
			case <-ct.refreshTicker.C:
				ct.refreshAll()
			}
		}
	}()
}
