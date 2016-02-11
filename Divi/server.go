package main 

import "net" 
import "fmt" 
import "io"

type client struct {
	message chan string
}

func main() {   
	fmt.Println("Starting Divi server")   
	// listen on all interfaces   
	ln, _ := net.Listen("tcp", ":8081")   
	// accept connection on port   

	// run loop forever (or until ctrl-c)   
	for { 
		conn, _ := ln.Accept()  

		cl := new(client)
		cl.message = make(chan string)
		go handleConnection(conn,cl.message)
                go handleMessage(cl.message)
	}
}

func handleMessage(msg chan string) {
    
   for {
	   select {
	   case m := <- msg: 
		   fmt.Print(m)                         		  
	   }
	   
   }
}



func handleConnection(c net.Conn,msg chan string) {

  

    buf := make([]byte, 4096)

    for {	   
	    
	    n, err := c.Read(buf)
            
             msg <- string(buf[0:n])
	    if err != nil || n == 0 {
		    c.Close()
		    break
	    }
	    n, err = c.Write(buf[0:n])
	    if err != nil {
		    c.Close()
		    break
	    }
    }
	
}

