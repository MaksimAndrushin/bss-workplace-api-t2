package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ozonmp/omp-demo-api/internal/app/retranslator"
)

func main() {

	ctx, cancelFunc := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)

	cfg := retranslator.Config{
		ChannelSize:   512,
		ConsumerCount: 2,
		ConsumeSize:   10,
		ProducerCount: 28,
		WorkerCount:   2,
		ConsumeTimeout: 5,
	}

	retranslator := retranslator.NewRetranslator(cfg)

	retranslator.Start(ctx)
	defer retranslator.Close()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	cancelFunc()

}
