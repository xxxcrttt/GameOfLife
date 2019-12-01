package main

import (
	"fmt"
	"strconv"
	"strings"
)

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p golParams, d distributorChans, alive chan []cell) {

	// Create the 2D slice to store the world.
	world := make([][]byte, p.imageHeight)
	for i := range world {
		world[i] = make([]byte, p.imageWidth)
	}

	// Request the io goroutine to read in the image with the given filename.
	d.io.command <- ioInput
	nameFile :=  strings.Join([]string{strconv.Itoa(p.imageWidth), strconv.Itoa(p.imageHeight)}, "x")
	d.io.filename <- nameFile

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
	// Calculate the new state of Game of Life after the given number of turns.
	for turns := 0; turns < p.turns; turns++ {

		//Created another 2D slice to store the world that has cache.
		tempWorld := make([][]byte, p.imageHeight)
		for i := range world {
			tempWorld[i] = make([]byte, p.imageWidth)
		}
		for y := 0; y < p.imageHeight; y++ {

			for x := 0; x < p.imageWidth; x++ {
				// Placeholder for the actual Game of Life logic: flips alive cells to dead and dead cells to alive.
				// Initialise the neighboursAlive to 0.
				var alive=0
				//If the cell is on the edge of the diagram, mod it to fix the rule of the game.
				//是一个球形的 计算机从0 -> 15, 所以要+/- 1
				var yy=y + p.imageHeight
				var xx=x + p.imageWidth
				if world [(yy-1)%p.imageHeight] [(xx-1)%p.imageWidth] ==255 {alive=alive+1}
				if world [yy%p.imageHeight]     [(xx-1)%p.imageWidth] ==255 {alive=alive+1}
				if world [(yy+1)%p.imageHeight] [(xx-1)%p.imageWidth] ==255 {alive=alive+1}

				if world [(yy-1)%p.imageHeight] [xx%p.imageWidth]     ==255 {alive=alive+1}
				if world [(yy+1)%p.imageHeight] [xx%p.imageWidth]     ==255 {alive=alive+1}

				if world [(yy-1)%p.imageHeight] [(xx+1)%p.imageWidth] ==255 {alive=alive+1}
				if world  [yy%p.imageHeight]    [(xx+1)%p.imageWidth] ==255 {alive=alive+1}
				if world [(yy+1)%p.imageHeight] [(xx+1)%p.imageWidth] ==255 {alive=alive+1}


				// When the colour is white, the cell status is alive, parameter is 255.
				// When the colour is black, the cell status is dead, parameter is 0.
				if world[y][x] == 255 {
					// The conditions of the game
					// If less than 2 or more than 3 neighbours, live cells dead.
					if alive<2 || alive > 3{
						tempWorld[y][x] = 0
					}
					if alive==2 || alive==3{
						tempWorld[y][x] = world[y][x]
					}
				}else if world[y][x] ==0 {
					// If 3 neighbours alive, dead cells alive.
					if alive==3{
						tempWorld[y][x] = 255
					}else{
						tempWorld[y][x]=0
					}
				}
			}
		}
		for y := 0; y < p.imageHeight; y++ {
			for x := 0; x < p.imageWidth; x++ {
				// Replace placeholder nextWorld[y][x] with the real world[y][x]
				world[y][x]=tempWorld[y][x]
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
	d.io.command <- ioOutput

	d.io.filename <- nameFile
	d.io.output <- world

	// Make sure that the Io has finished any output before exiting.
	d.io.command <- ioCheckIdle
	<-d.io.idle

	// Return the coordinates of cells that are still alive.
	alive <- finalAlive
}
