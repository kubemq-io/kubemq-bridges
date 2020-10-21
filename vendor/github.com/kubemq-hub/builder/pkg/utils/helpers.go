package utils

import (
	"bufio"
	"os"
)

func WaitForEnter() {
	Println("<cyan>Press ENTER to go back</>")
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
}
