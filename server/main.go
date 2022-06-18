package main

import (
	"context"
	"fmt"
	db "github/vikashparashar/gRPC_Project_01_Blog/database"
	pb "github/vikashparashar/gRPC_Project_01_Blog/proto"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BlogItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Author_id string             `bson:"author_id"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
}

func documentToBlog(in *BlogItem) *pb.Blog {
	return &pb.Blog{
		Id:       in.ID.Hex(),
		AuthorId: in.Author_id,
		Title:    in.Title,
		Content:  in.Content,
	}
}

var (
	network = "tcp"
	address = "0.0.0.0:50051"
	coll    *mongo.Collection
)

type Server struct {
	pb.Blog_ServiceServer
}

func main() {
	_, coll = db.StratDatabase()
	lis, err := net.Listen(network, address)
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()
	pb.RegisterBlog_ServiceServer(s, &Server{})
	if err = s.Serve(lis); err != nil {
		log.Println(err)
	}
}

func (s *Server) Create_One_Blog(ctx context.Context, in *pb.Blog) (*pb.BlogId, error) {
	data := BlogItem{
		Author_id: in.AuthorId,
		Title:     in.Title,
		Content:   in.Content,
	}
	log.Println("____Create_One_Blog Function Was Invoked At Server___")
	res, err := coll.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("internal error : %v\n", err),
		)
		// log.Fatalln("___failed to insert document into collection___")
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			"can not convert to oid",
		)
	}
	return &pb.BlogId{
		Id: oid.Hex(),
	}, nil

}

// func (s *Server) Create_Multiple_Blog(stream pb.Blog_Service_Create_Multiple_BlogServer) error {
// 	log.Println("____Create_Multiple_Blog Function Was Invoked At Server___")
// }
func (s *Server) Read_One_Blog(ctx context.Context, id *pb.BlogId) (*pb.Blog, error) {
	log.Println("____Read_One_Blog Function Was Invoked At Server___")
	var data *BlogItem
	oid, err := primitive.ObjectIDFromHex(id.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"can't parse ID",
		)
	}
	filter := bson.M{"_id": oid}
	res := coll.FindOne(context.Background(), filter)
	if err = res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"can't find the blog with given ID",
		)
	}
	blog := documentToBlog(data)
	return blog, nil

}

func (s *Server) Real_All_Blog(_ *emptypb.Empty, stream pb.Blog_Service_Real_All_BlogServer) error {
	log.Println("____Real_All_Blog Function Was Invoked At Server___")
	cur, err := coll.Find(context.Background(), primitive.D{})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("unknown internal error : %v\n", err),
		)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// stream.Send(cur.All(context.Background(), bson.D{}))
		data := &BlogItem{}
		err = cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error While Decoding Data From MongoDB : %v\n", err),
			)
		}
		if err = stream.Send(documentToBlog(data)); err != nil {
			log.Fatalf("Server Streaming Error : %v\n", err)
		}
	}
	if err = cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown Internal Error : %v\n", err),
		)
	}
	return nil
}
func (s *Server) Update_One_Blog(ctx context.Context, in *pb.Blog) (*emptypb.Empty, error) {
	log.Println("____Update_One_Blog Function Was Invoked At Server___")
	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"can't parse ID",
		)
	}
	data := BlogItem{
		Author_id: in.AuthorId,
		Title:     in.Title,
		Content:   in.Content,
	}
	filter := bson.M{"_id": oid}
	update := bson.M{"$set": data}
	res, err := coll.UpdateByID(context.Background(), filter, update)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"could not update",
		)
	}
	if res.MatchedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			"blog not found with given id",
		)
	}
	return &emptypb.Empty{}, nil

}

func (s *Server) Delete_One_Blog(ctx context.Context, id *pb.BlogId) (*emptypb.Empty, error) {
	log.Println("____Delete_One_Blog Function Was Invoked At Server___")
	oid, err := primitive.ObjectIDFromHex(id.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"can't parse ID",
		)
	}
	filter := bson.M{"_id": oid}

	res, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"could not delete",
		)
	}
	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			"blog not found with given id",
		)
	}
	return &emptypb.Empty{}, nil
}

// func (s *Server) Delete_All_Blog(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
// 	log.Println("____Delete_All_Blog Function Was Invoked At Server___")
// }
