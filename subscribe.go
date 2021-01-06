package forest

import (
	"bytes"
	"context"
	"time"
)

// SubscribeKeyValue listen to changes in KV store. the return value (a channel) of this method will receive new data and will be reused to send new data when there's changes of config in vault.
// SO DO NOT CLOSE THE CHANNEL LIKE EVER from your side unless you want some panic to happen.
//
// Because there is currently no support for socket connection for Vault (as of 6 Jan 2021), an approach mimicking concept `stale-while-revalidate` is used instead.
// Every `age` time passes, a request for new key value to vault is sent and compared with previous. If there's a difference, the new data is sent to the channel.
//
// Example:
//  chanConf := forest.SubscribeKeyValue(context.Background(), "some-conf", time.Second * 5, func(err) (exit bool) {
//    fmt.Println(err) // You get err when failed to get the resource. context cancel error (like timeout) is also send here.
//    // return true // setting return to true will make the program exit this SubscribeKeyValue method stack and the channel will be closed.
//    return false // Deliberately allows the function to continue
//  })
//  go func() { // Register changes to viper
//    viper.ReadConfig(bytes.NewBuffer(<-chanConf))
//  }()
func (v *Vault) SubscribeKeyValue(ctx context.Context, key string, age time.Duration, errFunc func(error) (exit bool)) <-chan []byte {
	oriCtx := ctx
	target := make(chan []byte)
	timeout := age
	if timeout < time.Second*5 {
		timeout = time.Second * 5
	}
	if errFunc == nil {
		errFunc = func(error) bool { return false }
	}
	go func() {
		defer close(target)
		ctx, clear := context.WithTimeout(ctx, timeout)
		prev, err := v.GetKeyValue(ctx, key)
		if err != nil {
			exit := errFunc(err)
			if exit {
				clear()
				return
			}
		}
		clear()
		for {
			select {
			case <-oriCtx.Done():
				return
			default:
				time.Sleep((age))
				ctx, clear := context.WithTimeout(ctx, timeout)
				next, err := v.GetKeyValue(ctx, key)
				if err != nil {
					exit := errFunc(err)
					if exit {
						clear()
						return
					}
					clear()
					continue
				}
				if !bytes.Equal(prev, next) {
					target <- next
					prev = next
				}
				clear()
			}
		}
	}()
	return target
}
