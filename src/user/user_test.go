package user

import (
  "testing"
  "io/ioutil"
  "os"
)

const user_db = "users.db"
const name = "testUser"
const pass = "testPass"

func TestNew(t *testing.T) {
  u := New(user_db, name)
  if u.userName != name {
    t.Errorf("user name not correct, expected %s and got %s", name, u.userName)
  }
}

func TestUserName(t *testing.T) {
  u := New(user_db, name)
  if u.UserName() != name {
    t.Errorf("user name not correct, expected %s and got %s", name, u.UserName())
  }
}

func TestAdd(t *testing.T) {
  u := New(user_db, name)
  err := u.Add(pass)
  if err != nil {
    t.Errorf("function returned error %s", err)
  }
  files, err := ioutil.ReadDir(".")
  found := false
  for _, file := range files {
    if file.Name() == user_db {
      found = true
    }
  }
  if !found {
    t.Errorf("The user database was not created.")
  }
  err = u.Add(pass)
  if err != USER_EXISTS {
    t.Errorf("function did not detect existsing user")
  }
  os.Remove(user_db)
}

func TestValidate(t *testing.T) {
  u := New(user_db, name)
  u.Add(pass)
  u2 := New(user_db, "testUser2")
  err := u2.Validate(pass)
  if err != USER_DOES_NOT_EXIST {
    t.Errorf("The function did not detect non existant user")
  }
  err = u.Validate("wrongPass")
  if err != INCORRECT_PASSWORD {
    t.Errorf("The function did not detect incorrect password.")
  }
  err = u.Validate(pass)
  if err != nil {
    t.Errorf("The function failed to authenticate with valid password.")
  }
  os.Remove(user_db)
}

func TestList(t *testing.T) {
  u := New(user_db, name)
  u.Add(pass)
  u2 := New(user_db, "testUser2")
  u2.Add(pass)
  users := List(user_db)
  if len(users) != 2 {
    t.Errorf("Incorect list length, expected 2 got %d", len(users))
  }
  if users[0] != name || users[1] != "testUser2" {
    t.Error("The users where not properly listed")
  }
  os.Remove(user_db)

}
