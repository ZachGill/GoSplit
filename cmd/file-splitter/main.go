package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "boom"
	app.Usage = "make an explosive entrance"
	app.Action = action

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "files, f",
			Value: 5,
			Usage: "Number of files per directory. 0 for no directory splitting",
		},
		cli.IntFlag{
			Name:  "lines, l",
			Value: 2e6,
			Usage: "Number of lines per file",
		},
		cli.StringFlag{
			Name:  "source, s",
			Value: "",
			Usage: "File to split",
		},
		cli.StringFlag{
			Name:  "destination, d",
			Value: "",
			Usage: "Destination to send split files to",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context) {
	var (
		err               error
		source            *os.File
		outputFile        *os.File
		fileName          string
		destination       string
		directory         string
		lineCount         int
		linesPerFile      int
		filesPerDirectory int
		lineCounter       int
		fileCounter       int
		directoryCounter  int
	)
	destination = c.String("destination")

	linesPerFile = c.Int("lines")

	log.Println("Opening source file")
	if source, err = os.Open(c.String("source")); err != nil {
		log.Println("Error opening file:", err)
	}
	defer source.Close()

	log.Println("Counting lines in file")
	if lineCount, err = countLines(source); err != nil {
		log.Println("Error counting lines:", lineCount)
	}
	source.Close()
	log.Println("line count:", lineCount)

	filesPerDirectory = c.Int("files")

	if source, err = os.Open(c.String("source")); err != nil {
		log.Println("Error opening file:", err)
	}
	scanner := bufio.NewScanner(source)
	lineCounter = 0
	fileCounter = 0
	directoryCounter = 0
	fileName = fmt.Sprintf("%s/delete_queries_%d%s", directory, directoryCounter, string(97+fileCounter))
	if outputFile, err = os.Create(fileName); err != nil {
		log.Println("error creating output file:", err.Error())
		//logger.Error("error creating output file", zap.Error(err))
		return
	}
	defer outputFile.Close()
	for scanner.Scan() {

		if fileCounter == 0 {
			directory = fmt.Sprintf("delete_queries_batch_%d", directoryCounter)
			log.Println(directoryCounter)
			if err = os.Mkdir(directory, 0777); err != nil {
				log.Println("error creating output directory:", err.Error())
				//logger.Error("error creating directory", zap.Error(err))
				break
			}
			log.Println("created new directory:", directory)
			directoryCounter++
		}

		if lineCounter == 0 {
			outputFile.Close()
			fileName = fmt.Sprintf("%s/%s/delete_queries_%d%s", destination, directory, directoryCounter, string(97+fileCounter))
			if outputFile, err = os.Create(fileName); err != nil {
				log.Println("error creating output file:", err.Error())
				//logger.Error("error creating output file", zap.Error(err))
				break
			}
			log.Println("created new file:", fileName)
			fileCounter++
		}

		outputString := fmt.Sprintf("%s\n", scanner.Text())
		if _, err = outputFile.WriteString(outputString); err != nil {
			log.Println("error writing line to file:", err.Error())
			break
		}

		lineCounter++

		if lineCounter == linesPerFile {
			lineCounter = 0
		}

		if fileCounter == filesPerDirectory {
			fileCounter = 0
		}
	}
	if err != nil {
		log.Println("unknown error:", err.Error())
		// logger.Error("unknown error", zap.Error(err))
	}
	log.Println("End")
}

func countLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
