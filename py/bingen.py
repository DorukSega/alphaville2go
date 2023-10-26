import random as rn

# Data to write to the binary file
w = 50
h = 50
data = b''
data += bytes([w])
data += bytes([h])

while len(data) < w*h:
    iy = rn.randint(0, 44)
    data += bytes([15])

# File path for the binary file
file_path = './tilemap.bin'

# Open the binary file in write binary mode
with open(file_path, 'wb') as file:
    file.write(data)

print(f"Binary data has been written to '{file_path}'.")
