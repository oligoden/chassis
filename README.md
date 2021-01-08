# chassis
A MVC that handles storage and authorisation.

device -> device/model -> storage
device -> device/model -> device/model/data -> storage
device -> device/view -> device/model
storage/gosql -> storage

## Installation

```bash
go get github.com/oligoden/chassis
```

## Usage

```golang
import "github.com/oligoden/chassis.git"

const (
	dbt = "mysql"
	uri = "chassis:password@tcp(localhost:3316)/chassis?charset=utf8&parseTime=True&loc=Local"
)

func main() {
    router := httprouter.New()
    
    // start storage
    store := gormdb.New(dbt, uri)

    //initialize device
	dProject := project.NewDevice(store)
    dProject.Manage("migrate")
    
    // use controllers on device as handlers
    router.Handler("POST", "/projects", dProject.Create())
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

## Testing

Spin up test database before running tests with:
 
```bash
docker-compose up -d
```