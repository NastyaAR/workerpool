package main

import (
	"bytes"
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	wPool, err := NewWorkerPool(ctx, 3, 6, 10)
	if err != nil {
		panic(err)
	}

	buf1 := bytes.NewBufferString("aaa")
	buf2 := bytes.NewBufferString("bbb")
	buf3 := bytes.NewBufferString("ccc")
	buf4 := bytes.NewBufferString("ddd")
	buf5 := bytes.NewBufferString("eee")
	buf6 := bytes.NewBufferString("fff")
	buf7 := bytes.NewBufferString("ggg")
	buf8 := bytes.NewBufferString("kkk")
	buf9 := bytes.NewBufferString("eeee")

	wPool.Submit(buf1)
	wPool.Submit(buf2)
	wPool.Submit(buf3)

	wPool.AddWorker(ctx)
	wPool.Submit(buf4)
	wPool.Submit(buf5)
	wPool.DeleteWorker()
	wPool.DeleteWorker()
	wPool.DeleteWorker()
	wPool.DeleteWorker()
	fmt.Println("delete all")

	wPool.Submit(buf6)
	wPool.Submit(buf7)

	wPool.AddWorker(ctx)
	wPool.AddWorker(ctx)

	wPool.Submit(buf8)
	wPool.Submit(buf9)

	wPool.DeleteWorker()
	wPool.DeleteWorker()

	fmt.Println("start wait")
	wPool.Wait()
	fmt.Println("stop wait")

	wPool.Free()

	wPool.AddWorker(ctx)
	wPool.AddWorker(ctx)

	buf8 = bytes.NewBufferString("kkk again")
	buf9 = bytes.NewBufferString("eeee again")

	wPool.Submit(buf8)
	wPool.Submit(buf9)

	wPool.DeleteWorker()
	wPool.DeleteWorker()

	wPool.Wait()

	fmt.Println("successfully returned")
}
