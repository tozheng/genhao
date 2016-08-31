package main

import (
    "io"
    "log"
    "net/http"
    "strconv"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "fmt"
    "time"
    "os"
    "bufio"
    "strings"
)

func readPWD(fileName string) string {
    f, err := os.Open(fileName)
    if err != nil {
        return ""
    }
    buf := bufio.NewReader(f)
    line, err := buf.ReadString('\n')
    line = strings.TrimSpace(line)
    return line
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    param1 := r.URL.Query().Get("name")
    param2 := r.URL.Query().Get("phone")

    pwd := readPWD("./genhao.config")
    fmt.Println(pwd)
    dbPath := "@tcp(genhao-instance.crfpanrt7dib.ap-northeast-2.rds.amazonaws.com:3306)/genhao?charset=utf8"
    fmt.Println(pwd+dbPath)
    db := opendb(pwd+dbPath)
    rows, err := db.Query("SELECT * FROM feed")

    if err == nil {
        var step int = 0

        io.WriteString(w, "{\"feed\": [")

        for rows.Next() {
            var fid int32
            var title string
            var detail string
            var uid int32
            var createTime int64
            var note string
            err = rows.Scan(&fid, &title, &detail, &uid, &createTime, &note)

            if step != 0 {
                io.WriteString(w, ",")
            }

            // title
            io.WriteString(w, "{\"title\":\"")
            io.WriteString(w, title)

            io.WriteString(w, "\",\"content\":\"")
            io.WriteString(w, detail)

            io.WriteString(w, "\",\"username\":\"")
            io.WriteString(w, queryUser(db, uid))

            io.WriteString(w, "\",\"time\":\"")
            io.WriteString(w, time.Unix(createTime, 0).Format("2006-01-02 15:04:05"))

            io.WriteString(w, "\",\"visitor\":\"")
            io.WriteString(w, param1)

            io.WriteString(w, "\",\"phone\":\"")
            io.WriteString(w, param2)
            io.WriteString(w, "\"")

            images := queryImage(db, fid)
            if len(images) > 0 {
                io.WriteString(w, ",\"imageName\":[")
                for index, v := range images {
                    if index > 0 {
                        io.WriteString(w, ",")
                    }

                    io.WriteString(w, "\"")
                    io.WriteString(w, v);
                    io.WriteString(w, "\"")
                }

                io.WriteString(w, "]")
            }

            io.WriteString(w, "}")

            step += 1
        }

        io.WriteString(w, "]}")
    }

    db.Close()
}

func main() {
    http.HandleFunc("/user", userHandler)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err.Error())
    }

//    db := opendb("tozheng:linkedin1@tcp(genhao-instance.crfpanrt7dib.ap-northeast-2.rds.amazonaws.com:3306)/genhao?charset=utf8")
//    userId := insertUser(db, "郑叶1", 38, "18800000001", "haha")
//    feedId := insertFeed(db, "Feed1", "Feed content 1", userId, "feed note 1")
//    insertImage(db, feedId, "feed image", "http://www.baidu.com/path/1111.jpg")
//    insertImage(db, feedId, "feed image", "http://www.baidu.com/path/2222.jpg")
}

//打开数据库连接
func opendb(dbstr string) ( * sql.DB) {
    //dsn: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&paramN=valueN]
    db, err := sql.Open("mysql", dbstr)
    if err != nil {
        log.Fatalf("Open database error: %s\n", err)
    }

    return db
}

//插入user数据
func insertUser(db  * sql.DB, name string, age int32, phone string, note string) int32 {
    stmt, err := db.Prepare("INSERT INTO user SET name=?,age=?,phone=?,note=?")
    defer stmt.Close()

    if err != nil {
        log.Println(err)
        return 0
    }

    res, _ := stmt.Exec(name, age, phone, note)
    id, _ := res.LastInsertId()

    fmt.Println(id)
    return int32(id)
}

//插入feed数据
func insertFeed(db  * sql.DB, title string, detail string, userId int32, note string) int32 {
    stmt, _ := db.Prepare("INSERT INTO feed SET title=?,detail=?,user=?,time=?,note=?")
    defer stmt.Close()

    res, _ := stmt.Exec(title, detail, userId, time.Now().Unix(), note)
    id, _ := res.LastInsertId()

    fmt.Println(id)
    return int32(id)
}

//插入image数据
func insertImage(db  * sql.DB, fid int32, note string, path string) int32 {
    stmt, _ := db.Prepare("INSERT INTO image SET fid=?,time=?,note=?,path=?")
    defer stmt.Close()

    res, _ := stmt.Exec(fid, time.Now().Unix(), note, path)
    id, _ := res.LastInsertId()

    fmt.Println(id)
    return int32(id)
}

func queryImage(db  * sql.DB, feedId int32) (images []string)  {
    feedimages := make([]string, 0, 10)

    queryStr1 := "SELECT * FROM image where fid="
    queryStr2 := strconv.Itoa(int(feedId))
    rows, err := db.Query(queryStr1+queryStr2)

    fmt.Println(queryStr1+queryStr2)

    for rows.Next() {
        var iid int32
        var fid int32
        var createTime int64
        var note string
        var path string
        err = rows.Scan(&iid, &fid, &createTime, &note, &path)
        if (err == nil) {
            feedimages = append(feedimages, path)
        }
    }

    return feedimages
}

func queryUser(db  * sql.DB, userId int32) (string)  {
    queryStr1 := "SELECT * FROM user where uid="
    queryStr2 := strconv.Itoa(int(userId))
    rows, err := db.Query(queryStr1+queryStr2)

    fmt.Println(queryStr1+queryStr2)

    for rows.Next() {
        var uid int32
        var name string
        var age int32
        var phone string
        var note string
        err = rows.Scan(&uid, &name, &age, &phone, &note)
        if (err == nil) {
            return name
        }
    }

    return ""
}
