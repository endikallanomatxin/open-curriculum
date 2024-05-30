# Open curriculum

Made with go and postgresql.

## Manual

To run it:

```bash
sudo docker compose up --build 
```

```bash
sudo docker compose -f docker-compose.dev.yml up --build
```

To use the db:

```bash
psql mydatabase -U myusername
```

To see all tables:
    
```sql
\dt
```

To see all rows in a table:

```sql
SELECT * FROM mytable;
```


## Resources

Air for live reload:
https://medium.easyread.co/today-i-learned-golang-live-reload-for-development-using-docker-compose-air-ecc688ee076

SSL certificate:
https://goenning.net/blog/free-and-automated-ssl-certificates-with-go/
