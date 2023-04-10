
# Response Time Simulation Service

## Deskripsi

Service ini bertujuan untuk mensimulasikan waktu respons dari request dengan jumlah data tertentu.

## Setup

1. Pastikan Anda telah memiliki database Postgre SQL yang terinstall di komputer Anda.
2. Buatlah database baru dengan nama "simulation" pada database Postgre SQL dengan menggunakan perintah SQL berikut:

> CREATE DATABASE simulation;

4. Setelah itu, clone repository project dengan menggunakan command:

> git clone https://github.com/sandisuryadi36/response-time-simulation.git

 4. Masuk ke dalam folder project dengan menggunakan command: 

> cd response-time-simulation

 5. Install semua dependencies yang dibutuhkan dengan menggunakan command:
 6. Sesuaikan file `.env` dengan configurasi Postgre SQL anda.

> go mod download


 ``## Menjalankan Service

1. Setelah melakukan setup, jalankan command `go run ./server` pada directory root project untuk memulai server.

 2. Server akan berjalan pada `localhost:8080`.
3. Lakukan request pada url `localhost:8080/api/store` dengan menggunakan aplikasi untuk melakukan request HTTP, seperti Postman atau Insomnia.
4. Pastikan struktur data JSON yang dikirimkan sesuai dengan format yang telah ditentukan, seperti contoh berikut:

>     { 
>         "request_id": 1234555, 
>         "data": [ 
>     	    { 
>     		    "id": 1234, 
>     		    "customer": "Jhon Smith", 
>     		    "quantity": 1, 
>     		    "price": 10.00, 
>     		    "timestamp": "2022-01-02 22:10:44" 
>     	    } 
>     	  ] 
>     }

 5. Setelah request berhasil terkirim, Anda akan menerima response dengan response time dari request tersebut.`
