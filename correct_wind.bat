@echo off
nircmd exec show "HD-Player" --instance Nougat32
nircmd wait 3000
nircmd win setsize title "BlueStacks" 0 0 600 1000
exit