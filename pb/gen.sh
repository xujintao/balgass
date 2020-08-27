export PROJECT_HOME=~/github.com/xujintao/balgass

# game_data
protoc \
--proto_path=$PROJECT_HOME/pb \
--go_out=$PROJECT_HOME/cmd/server_data/model \
game_data.proto