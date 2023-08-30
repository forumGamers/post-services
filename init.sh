docker exec -it database mongosh

rs.initiate(
    { 
        _id:"db-dev",
        members:[ 
            { _id:0 ,host:"database:27017"}, 
            { _id:1 , host:"database-secondary:27017" },
            { _id:2 , host:"database-arbiter:27017",arbiterOnly:true} 
            ]
    }
)