import requests
import json
import time

def get_block_size(block_height):
    url = f"https://lcd.archive.osmosis.zone/cosmos/base/tendermint/v1beta1/blocks/{block_height}"
    response = requests.get(url)
    if response.status_code == 200:
        block_data = response.json()
        
        # Remove the `sdk_block` key from the response, if it exists
        if 'sdk_block' in block_data:
            del block_data['sdk_block']
        
        # Convert the remaining JSON data back to a string to calculate size
        block_json_str = json.dumps(block_data)
        
        return len(block_json_str.encode('utf-8')) / (1024 * 1024)  # Convert bytes to megabytes
    else:
        print(f"Failed to retrieve block {block_height}: Status code {response.status_code}")
        return 0

def get_total_and_average_block_size(start_block, end_block):
    total_size_mb = 0
    block_count = 0
    
    for block_height in range(start_block, end_block + 1):
        block_size_mb = get_block_size(block_height)
        total_size_mb += block_size_mb
        block_count += 1
        if block_size_mb > 2:
            print(f"Block {block_height}: Size {block_size_mb:.6f} MB")
        
        # Introduce a 1-second delay every 10 blocks
        if block_count % 10 == 0:
            time.sleep(2)
    
    average_size_mb = total_size_mb / block_count if block_count > 0 else 0
    return total_size_mb, average_size_mb

# Calculate the block range for the last two weeks
latest_block = 19284772
blocks_in_two_weeks = 604800
start_block = latest_block - blocks_in_two_weeks
end_block = latest_block

# Calculate the total and average size for blocks in the last two weeks
total_size_mb, average_size_mb = get_total_and_average_block_size(start_block, end_block)

print(f"Total size from block {start_block} to {end_block}: {total_size_mb:.6f} MB")
print(f"Average block size from block {start_block} to {end_block}: {average_size_mb:.6f} MB")


