# VCS_SMS

1. Khởi tạo file app.env giống trong file app.txt
2. Khởi tạo container docker :

- Dùng lệnh docker-compose up -d

3. Kết nối Database (PostgresSQL / PgAdmin4)"

- Run file migrate.go
- Mở PgAdmin4, tạo một server mới với các thông tin giống trong app.env :
  - home/address: localhost
  - Port: 6500
  - Maintainance DB: VCS_SMS
  - Username: Postgres

4. Chạy chương trình:

- Run lệnh "air"

5. Sử dụng Postman để test:

![image](https://user-images.githubusercontent.com/88451173/196332348-15dab904-3048-4e3a-bb6a-17fd8fc734c6.png)

- Register:

![image](https://user-images.githubusercontent.com/88451173/196332443-ce5d4fcd-34c1-4c88-b4cd-384c3d5d24ca.png)

- Login:

![image](https://user-images.githubusercontent.com/88451173/196332528-ceb76288-7736-4c2b-9ea4-15116a2bfa32.png)

Sau khi login hệ thống sẽ trả lại một access token cho người dùng để thực hiện các quyền tiếp theo

- Check status of IpV4:

Hệ thống rà 1 lượt và ping đến các địa chỉ Ipv4 có trong DB, trả lại updated status của từng server

![image](https://user-images.githubusercontent.com/88451173/196333916-97413fb6-209c-48f5-aa50-f62b487035e4.png)

- Add servers:

![image](https://user-images.githubusercontent.com/88451173/196332730-6f156ea5-6270-480c-bc1e-2a4187664003.png)

- Change servers:

![image](https://user-images.githubusercontent.com/88451173/196332833-80f5ed34-3354-4d82-b8bb-1b0d588100fd.png)

- Get servers (with filter from/to)

![image](https://user-images.githubusercontent.com/88451173/196332943-339b263e-2d0b-4ea3-b439-1ac5446f3c3a.png)

- Delete server:

![image](https://user-images.githubusercontent.com/88451173/196333019-4df4fc38-4785-405b-be2c-68d652e65f6a.png)

- Import data from file excel:

![image](https://user-images.githubusercontent.com/88451173/196333117-bc44987c-8edd-4033-bf6e-c2db96bb611d.png)

Với file excel có dạng như sau:

![image](https://user-images.githubusercontent.com/88451173/196333228-cd80fd62-650b-4361-90a3-bc6d2cc9d96f.png)

Cột 1: Server_id,
Cột 2: Server_name,
Cột 3: Status,
Cột 4: CurrentUser_id,
Cột 5: Created_time,
Cột 6: Last_updated_time,
Cột 7: IpV4.

- Export file excel:
  Lấy mọi data trong DB chuyển thành file excel:

![image](https://user-images.githubusercontent.com/88451173/196333631-327f9693-0947-400d-bf0a-6846d4f9e33d.png)

- Chưa hoàn thành:
  UnitTest (đã hiểu và đang làm lại)
  Realtime Update
  calculate Uptime using Elasticsearch
