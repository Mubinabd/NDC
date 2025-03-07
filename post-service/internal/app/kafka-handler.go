package app

import (
	"context"
	"log"

	"google.golang.org/protobuf/encoding/protojson"
	pb "posts/internal/pkg/genproto"
	"posts/internal/usecase/service"
)

type KafkaHandler struct {
	log  *service.LogService
	post *service.PostService
}

func (h *KafkaHandler) LogUpdate() func(message []byte) {
	return func(message []byte) {

		//unmarshal the message
		var cer pb.LogUpdateRequest
		if err := protojson.Unmarshal(message, &cer); err != nil {
			log.Fatalf("Failed to unmarshal JSON to Protobuf message: %v", err)
			return
		}

		res, err := h.log.Update(context.Background(), &cer)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("Updated Log: %+v", res)
	}
}

func (h *KafkaHandler) LogDelete() func(message []byte) {
	return func(message []byte) {

		//unmarshal the message
		var cer pb.GetId
		if err := protojson.Unmarshal(message, &cer); err != nil {
			log.Fatalf("Failed to unmarshal JSON to Protobuf message: %v", err)
			return
		}

		res, err := h.log.Delete(context.Background(), &cer)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("Deleted Log: %+v", res)
	}
}

func (h *KafkaHandler) LogCreate() func(message []byte) {
	return func(message []byte) {

		//unmarshal the message
		var cer pb.LogCreateRequest
		if err := protojson.Unmarshal(message, &cer); err != nil {
			log.Fatalf("Failed to unmarshal JSON to Protobuf message: %v", err)
			return
		}

		res, err := h.log.Create(context.Background(), &cer)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("Created Log: %+v", res)
	}
}

func (h *KafkaHandler) PostUpdate() func(message []byte) {
	return func(message []byte) {

		//unmarshal the message
		var cer pb.PostUpdateRequest
		if err := protojson.Unmarshal(message, &cer); err != nil {
			log.Fatalf("Failed to unmarshal JSON to Protobuf message: %v", err)
			return
		}

		res, err := h.post.Update(context.Background(), &cer)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("Updated Post: %+v", res)
	}
}

func (h *KafkaHandler) PostDelete() func(message []byte) {
	return func(message []byte) {

		//unmarshal the message
		var cer pb.GetById
		if err := protojson.Unmarshal(message, &cer); err != nil {
			log.Fatalf("Failed to unmarshal JSON to Protobuf message: %v", err)
			return
		}

		res, err := h.post.Delete(context.Background(), &cer)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("Deleted Post: %+v", res)
	}
}

func (h *KafkaHandler) PostCreate() func(message []byte) {
	return func(message []byte) {

		//unmarshal the message
		var cer pb.PostCreateRequest
		if err := protojson.Unmarshal(message, &cer); err != nil {
			log.Fatalf("Failed to unmarshal JSON to Protobuf message: %v", err)
			return
		}

		res, err := h.post.Create(context.Background(), &cer)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("Created Post: %+v", res)
	}
}
