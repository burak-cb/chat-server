package Others

/*func main() {
	// Reading a File.

	inputFile, fileOpeningError := os.Open("hitchhikersGuide.txt")

	if fileOpeningError != nil {
		log.Fatalln("ERROR:", fileOpeningError)
	}

	defer func(inputFile *os.File) {
		fileClosureError := inputFile.Close()
		if fileClosureError != nil {

		}
	}(inputFile)

	fileScanner := bufio.NewScanner(inputFile)

	for fileScanner.Scan() {
		fmt.Println(fileScanner.Text())
		time.Sleep(time.Second)
	}
}*/
