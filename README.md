### Simple Singleton MongoDB Session

#### Usages

##### Before anything else you should initail the connection ( session )
```go

func main() {
    ...
    err := mongodb.InitMongoDB("mongodb://localhost:27017")
    //It will initail the singleton of mongo session it will return error if you want to handle it
    ...
}
```

```go
type Document struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	Data      string        `json:"data" bson:"data"`
}

}

func ExecuteSomeStuffInMongoDB() {
   doc := Document{
		CreatedAt: time.Now(),
		Data:      "My Document",
	}
    col := mongodb.GetMgoCollection("mydb","mycollection")
    if err := col.Insert(&doc); err != nil {
        //Handle error if you want
    }
}

```