创建执行命令
go build -o bc main.go

./bc
Usage:
  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS
  printchain - Print all the blocks of the blockchain
  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO address.

采用bolt数据库： 以hash为key，block序列化为value存储



