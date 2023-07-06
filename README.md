##### README #####
##### NOTES #####
1. Payment types:-
    0 - Free
    1 - Cash
    2 - QR Maybank
    3 - QR TnG
2. start_item_bal and end_item_bal data type is string so that I can format the data as `1=150&2=45`, 1 indicating item_id 1 and 2 indicating item_id 2

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
