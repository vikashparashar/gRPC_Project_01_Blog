package main

import (
	"context"
	"fmt"
	pb "github/vikashparashar/gRPC_Project_01_Blog/proto"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewBlog_ServiceClient(conn)
	id := Create_One_Blog(c) // create a single blog mongoDB
	// Create_Multiple_Blog(c)          // create multiple blog in mongoDB
	Read_One_Blog(c, id)               // with vaild id
	Read_One_Blog(c, "fffa5af5689fa6") // with invalid id
	Update_One_Blog(c, id)             // update a single blog in mongoDB
	Real_All_Blog(c)                   // find all blogs from mongodb
	Delete_One_Blog(c, id)             // delete a single blogsss from mongoDB
	// Delete_All_Blog(c)
}

func Create_One_Blog(c pb.Blog_ServiceClient) string {
	log.Println("____Create_One_Blog Function Was Invoked At Client___")
	var req = &pb.Blog{
		AuthorId: "Vikash Parashar",
		Title:    "Hello",
		Content:  "Content Of First Blog : Golang gRPC Course",
	}
	res, err := c.Create_One_Blog(context.Background(), req)
	if err != nil {
		log.Fatalf("___Unexcepted error  : %v\n", err)
	}
	fmt.Printf("___Document Inserted Successfully With InsertedID: %v\n", res.Id)
	return res.Id
}

// func Create_Multiple_Blog(c pb.Blog_ServiceClient) {
// 	log.Println("____Blog Function Was Invoked At Client___")
// }

func Read_One_Blog(c pb.Blog_ServiceClient, id string) *pb.Blog {

	log.Println("____Read_One_Blog Function Was Invoked At Client___")
	// var data *pb.Blog
	// oid, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	req := &pb.BlogId{Id: id}
	res, err := c.Read_One_Blog(context.Background(), req)
	if err != nil {
		log.Printf("Error happned while reading : %v\n", err)
	}
	log.Printf("Blog was read : %v\n", res)
	return res
}

func Real_All_Blog(c pb.Blog_ServiceClient) {
	log.Println("____Real_All_Blog Function Was Invoked At Client___")
	stream, err := c.Real_All_Blog(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("something happned : %v\n", err)
		}
		fmt.Printf("Real_All_Blog Result : %v\n", res)
	}
}

func Update_One_Blog(c pb.Blog_ServiceClient, id string) {
	log.Println("____Update_One_Blog Function Was Invoked At Client___")
	newblog := &pb.Blog{
		Id:       id,
		AuthorId: "Khushboo Parashar",
		Title:    "A Beautiful Day",
		Content:  "Some Amazing Content",
	}
	_, err := c.Update_One_Blog(context.Background(), newblog)
	if err != nil {
		log.Fatalf("Error happened while updateing : %v\n", err)
	}
	log.Println("Blog was updated !")

}

func Delete_One_Blog(c pb.Blog_ServiceClient, id string) {
	log.Println("____Delete_One_Blog Function Was Invoked At Client___")
	req := &pb.BlogId{Id: id}
	_, err := c.Delete_One_Blog(context.Background(), req)
	if err != nil {
		log.Fatalf("Error While Deleting Blog : %v\n", err)
	}
	log.Println("Blog was deleted !")
}

// func Delete_All_Blog(c pb.Blog_ServiceClient) {
// 	log.Println("____Blog Function Was Invoked At Client___")
// }
