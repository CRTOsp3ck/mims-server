##### README #####
1. Install go
2. Install postgres
3. add root as psql user
4. set root psql password
5. create db with name root
6. git clone mims-server
7. change main.go db settings
8. change database.yml settings
9. i know it is shit systems design atm, but its a start

##### NOTES #####
1. Payment types:-
    1 - Cash
    2 - QR Maybank
    3 - QR TnG
    99 - Free
2. Item id:-
    1 - MD2 Juice
    2 - MD2 Sliced Fruit
3. start_item_bal and end_item_bal data type is string so that I can format the data as `1=150&2=45`, 1 indicating item_id 1 and 2 indicating item_id 2


##### 5/7/2023 #####
1. Create GET, POST for all tables:-
    # Small tables
    Agent
    Inventory
    Balances
    Item
    # Big tables
    Operation
    Sales

2. Remove item_id from inventory table and change 
