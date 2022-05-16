#!/usr/bin/env python
# coding: utf-8

import hashlib
import time

import bcrypt
import gevent
import pymysql


class MyPyMysql:
    def __init__(self, host, port, username, password, db, charset="utf8"):
        self.host = host
        self.port = port
        self.username = username
        self.password = password
        self.db = db
        self.charset = charset
        self.pymysql_connect()

    def pymysql_connect(self):
        self.conn = pymysql.connect(
            host=self.host,
            port=self.port,
            user=self.username,
            password=self.password,
            db=self.db,
            charset=self.charset,
        )

        self.asynchronous()

    def run(self, nmin, nmax):
        self.cur = self.conn.cursor()

        sql = "INSERT INTO users(username, password, profile_pic) VALUES (%s, %s, %s)"

        data_list = []
        # hashed = bcrypt.hashpw("test".encode("utf-8"), bcrypt.gensalt())
        hashed = hashlib.md5("test".encode())
        in_sql = hashed.hexdigest()
        for i in range(nmin, nmax):
            username = "test" + str(i)
            result = (username, in_sql, "./assets/default.jpeg")
            data_list.append(result)

        content = self.cur.executemany(sql, data_list)
        if content:
            print("成功插入第{}条数据".format(nmax - 1))

        self.conn.commit()

    def asynchronous(self):
        max_line = 10000
        g_l = [
            gevent.spawn(self.run, i, i + max_line)
            for i in range(1, 10000000, max_line)
        ]

        gevent.joinall(g_l)
        self.cur.close()
        self.conn.close()


if __name__ == "__main__":
    start_time = time.time()
    st = MyPyMysql("localhost", 3306, "root", "CdB5f2vY", "user")
    print("程序耗时{:.2f}".format(time.time() - start_time))
