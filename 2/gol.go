package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)
//判断alive的邻居
func numNeiAlive(world [][]byte,y,x int, imageHeight int,imageWidth int) int {
	var alive = 0
	x = x + imageWidth
	if world [y-1] [(x-1)%imageWidth] ==255 {alive=alive+1}
	if world [y]   [(x-1)%imageWidth] ==255 {alive=alive+1}
	if world [y+1] [(x-1)%imageWidth] ==255 {alive=alive+1}

	if world [y-1] [x%imageWidth]     ==255 {alive=alive+1}
	if world [y+1] [x%imageWidth]     ==255 {alive=alive+1}

	if world [y-1] [(x+1)%imageWidth] ==255 {alive=alive+1}
	if world [y]   [(x+1)%imageWidth] ==255 {alive=alive+1}
	if world [y+1] [(x+1)%imageWidth] ==255 {alive=alive+1}
	return alive
}

//slices
func buildNewWorld (world [][]byte, heightOfWorker, imageHeight, imageWidth, totalThreads, currentThreads int) [][] byte{
	newWorld := make([][]byte, heightOfWorker+2)
	for j := 0;j<heightOfWorker+2; j++ {
		newWorld[j] = make([]byte, imageWidth)
	}

	if currentThreads==0{
		for x := 0; x < imageWidth; x++ {
			newWorld[0][x]=world[imageHeight-1][x]
		}
	}else{
		for x := 0; x < imageWidth; x++ {
			newWorld[0][x]=world[currentThreads*heightOfWorker-1][x]
		}
	}

	for y := 1; y <= heightOfWorker; y++ {
		for x := 0; x < imageWidth; x++ {
			newWorld[y][x]=world[currentThreads*heightOfWorker+y-1][x]
		}
	}

	if currentThreads==totalThreads-1{
		for x := 0; x < imageWidth; x++ {
			newWorld[heightOfWorker+1][x]=world[0][x]
		}
	}else {
		for x := 0; x < imageWidth; x++ {
			newWorld[heightOfWorker+1][x]=world[(currentThreads+1)*heightOfWorker][x]
		}
	}
	return newWorld
}

//改变memory sharing -- worker function
func worker(wChan chan byte, imageHeight int,imageWidth int,  out chan byte){
	//整合 打包
	//create a new world with the channel
	world := make([][]byte, imageHeight + 2)
	for i := range world {
		world[i] = make([]byte, imageWidth)
	}
	for y := 0; y < imageHeight;y++{
		for x := 0; x < imageWidth; x++{
			world[y][x] =<- wChan
		}
	}

	tempWorld := make([][]byte, imageHeight + 2)
	for i := range world {
		tempWorld[i] = make([]byte, imageWidth)
	}

	for y := 1; y <= imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			// Placeholder for the actual Game of Life logic: flips alive cells to dead and dead cells to alive.
			var alive=0
			alive=numNeiAlive(world,y,x,imageHeight,imageWidth)

			//conditions of the game
			if world[y][x] == 255 {
				if alive<2 {
					tempWorld[y][x] = 0
				}
				if alive==2 || alive==3 {
					tempWorld[y][x] = world[y][x]
				}
				if alive>3 {
					tempWorld[y][x] = 0
				}
			}else if world[y][x] ==0 {
				if alive==3 {
					tempWorld[y][x] = 255
				} else {
					tempWorld[y][x] = 0
				}
			}

		}
	}
	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			out <- tempWorld[y+1][x]
		}
	}
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p golParams, d distributorChans, alive chan []cell, key chan rune) {

	// Create the 2D slice to store the world.
	world := make([][]byte, p.imageHeight)
	for i := range world {
		world[i] = make([]byte, p.imageWidth)
	}

	// Request the io goroutine to read in the image with the given filename.
	d.io.command <- ioInput
	d.io.filename <- strings.Join([]string{strconv.Itoa(p.imageWidth), strconv.Itoa(p.imageHeight)}, "x")
	//d.io.filename <- nameFile

	// The io goroutine sends the requested image byte by byte, in rows.
	for y := 0; y < p.imageHeight; y++ {
		for x := 0; x < p.imageWidth; x++ {
			val := <-d.io.inputVal
			if val != 0 {
				fmt.Println("Alive cell at", x, y)
				world[y][x] = val
			}
		}
	}

	// Create a 2s timer called ticker.
	ticker := time.NewTicker(2 * time.Second)
	out := make([]chan byte, p.threads)

	// Calculate the new state of Game of Life after the given number of turns.
	for turns := 0; turns < p.turns; turns++ {
		heightOfWorker := p.imageHeight / p.threads

		//running := true
		select {
		//case <-ticker.C:
		case press := <- key:
			if press == 's' {
				dio(d, p, world,turns)
//				d.io.command <- ioOutput
//				d.io.filename <- strings.Join([]string{strconv.Itoa(p.imageWidth), strconv.Itoa(p.imageHeight), strconv.Itoa(turns)}, "x")
//				d.io.output <- world
			}else if press == 'q' {
				dio(d, p, world,turns)
				//d.io.command <- ioOutput
				//d.io.filename <- strings.Join([]string{strconv.Itoa(p.imageWidth), strconv.Itoa(p.imageHeight), strconv.Itoa(turns)}, "x")
				//d.io.output <- world
				fmt.Println("terminate...")
				return
			} else if press =='p' {
				fmt.Println(turns)
				for {
						pp := <-key
						if pp == 'p' {
							fmt.Println("continuing...")
							break
						}
					}
				}

		case <- ticker.C:
			var finalAlive []cell
			// Go through the world and append the cells that are still alive.
			for y := 0; y < p.imageHeight; y++ {
				for x := 0; x < p.imageWidth; x++ {
					if world[y][x] != 0 {
						finalAlive = append(finalAlive, cell{x: x, y: y})
					}
				}
			}
			fmt.Println("number of alive cells is:", len(finalAlive))

		default:
		}

		//Put logic outside of select such that when other cases run, the logic doesn't skip a turn.

		//var out [8] chan [][]byte
	//IssueIssueIssueIssueIssueIssueIssue: 16 must be a constant number, it might goes wrong when it goes bigger.

	for i := 0; i < p.threads; i++ {

			out[i] = make(chan byte)
			wChan := make(chan byte)
			//build slices the workers need to work on
			newWorld := buildNewWorld(world,  heightOfWorker, p.imageHeight, p.imageWidth, p.threads, i)
			go worker(wChan, heightOfWorker+2 ,p.imageWidth , out[i])
			//Send world cells to workers
			for y := 0; y < heightOfWorker+2; y++ {
				for x := 0; x < p.imageWidth; x++ {
					wChan <- newWorld[y][x]
				}
			}
		}
		for i:=0; i<p.threads ; i++ {
			//slices from workers
			tempOut := make([][]byte, heightOfWorker)
			for i := range tempOut {
				tempOut[i] = make([]byte, p.imageWidth)
			}
			for y := 0; y < heightOfWorker; y++ {
				for x := 0; x < p.imageWidth; x++ {
					tempOut[y][x] = <-out[i]
				}
			}
			//println("tempOut  i=",i)
			for y := 0; y < heightOfWorker; y++ {
				for x := 0; x < p.imageWidth; x++ {
					//print(tempOut[y+1][x])
					world[i*heightOfWorker+y][x]=tempOut[y][x]
				}
				//println()
			}
		}

	}

	// Create an empty slice to store coordinates of cells that are still alive after p.turns are done.
	var finalAlive []cell
	// Go through the world and append the cells that are still alive.
	for y := 0; y < p.imageHeight; y++ {
		for x := 0; x < p.imageWidth; x++ {
			if world[y][x] != 0 {
				finalAlive = append(finalAlive, cell{x: x, y: y})
			}
		}
	}



	// Request the io goroutine to write in the image with the given filename.
	//d.io.command <- ioOutput
	//d.io.filename <- strings.Join([]string{strconv.Itoa(p.imageWidth), strconv.Itoa(p.imageHeight), strconv.Itoa(p.turns)}, "x")

	// Send the world to finalBoard
	//d.io.output <- world

	// Make sure that the Io has finished any output before exiting.
	d.io.command <- ioCheckIdle
	<-d.io.idle

	// Return the coordinates of cells that are still alive.
	alive <- finalAlive
}

func dio(d distributorChans, p golParams, world[][]byte, turn int) {
	d.io.command <- ioOutput
	d.io.filename <- strings.Join([]string{strconv.Itoa(p.imageWidth), strconv.Itoa(p.imageHeight), strconv.Itoa(turn)}, "x")
	d.io.output <- world
}
