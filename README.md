# go-gen-api

An attempt to auto generate an http rest api for a database without much fiddeling.

## What is that for?
Say you have following tables in a database and want to add an REST Api for it:

    [ User ]      [ Group ]    [ GroupMember ]
    - ID          - ID         - ID
    - Name        - Name       - UserID
    - Password                 - GroupID


The normal way would be to create all kind of CRUD operations for each Table:

    func CreateUser(w http.ResponseWriter, r *http.Request)  {
        if r.Method != "POST" || r.Method != "PUT" {
            throw some error
        }
        if err := json.unmarshal(r.Body, &newUserStruct); err != nil {
            throw some error
        }
        if err := validateUserStruct(newUserStruct); err != nil {
            throw some error
        }
        if err := createNewUserOnTable(&newUserStruct); err != nil {
            throw some error
        }
        send everything worked response
    }

    func UpdateUser(w http.ResponseWriter, r *http.Request)  {
        ...
    }

    func GetUser(w http.ResponseWriter, r *http.Request)  {
        ...
    }
    func DeleteUser(w http.ResponseWriter, r *http.Request)  {
        ...
    }
    ...

This is not only boring, it is also time consuming.  
With **go-gen-api** you can automate this process:

    type User struct {
        ID       int
        Name     string
        Password string
    }

    type Group struct {
        ID       int
        Name     string
    }

    type GroupMember struct {
        ID       int
        UserID   int
        GroupID  int
    }

    err := gogenapi.Generate(&gogenapi.Config{
		Structs:      []interface{}{&User{}, &Group{}, &GroupMember{}},
		OutputPath:   "generated",
	})
	if err != nil {
		panic(err)
	}

After generation you can use it in your project:

    db, err := sql.Open("sqlite3", "user.sqlite")
    router := mux.NewRouter()
	restAPI := generated.NewUserRestAPI(router.PathPrefix("/user").Subrouter(), generated.NewUserAPI(db))

    // You could also setup some hooks before or after execution
	restAPI.Hooks.PreCreate = func(r *http.Request, user *generated.User) error {
		if user.Name == nil || len(*user.Name) <= 0 {
			return errors.New("Invalid Name")
		}
		if user.Password == nil || len(*user.Password) <= 0 {
			return errors.New("Invalid Password")
		}
		return nil
	}

You can find a working example [here](https://github.com/Eun/loginexample)


### Limitations
- **go-gen-db does not create any tables**, you must create your tables before
- Any other fields types than string or int, (pull requests are welcome)

### What needs to be done
- [ ] Write some kind of tests
- [ ] Parse struct tags to ignore fields
- [ ] Parse and store other data than int and strings

