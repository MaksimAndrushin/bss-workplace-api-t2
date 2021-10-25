package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ozonmp/omp-demo-api/internal/app/retranslator"
)

func main() {

	var ctx = context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)

	var sigs = make(chan os.Signal, 1)

	var cfg = retranslator.Config{
		ChannelSize:   512,
		ConsumerCount: 2,
		ConsumeSize:   10,
		ProducerCount: 28,
		WorkerCount:   2,
		ConsumeTimeout: 5,
		Ctx: ctx,
		CancelFunc: cancelFunc,
	}

	var retranslator = retranslator.NewRetranslator(cfg)

	retranslator.Start()
	defer retranslator.Close()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

}
