### Simple Singleton MongoDB Session

#### Usages

##### Before anything else you should initail the connection ( session )
```go

import "github.com/liontail/mongodb"

func main() {
    ...
    sess, err := mongodb.InitMongoDB("mongodb://localhost:27017")
    //It will initail the singleton of mongo session it will return error if you want to handle it
    defer sess.close()
    // after end process main will close the session

    sess.Options.ReconnectTime = 10 // default 15
    // there are keep alive session, overwrite this option to change the time or remove this to use default
    ...
}
```

```go
import "github.com/liontail/mongodb"

type Document struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	Data      string        `json:"data" bson:"data"`
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

    id := "34yhgvfrtyui213" // id is the hex of the objectId
    findDoc := Document{}
    if err := col.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&findDoc); err != nil {
        //Handle error if you want
    }
}

```
