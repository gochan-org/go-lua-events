print("Testing post-received event in eventstest.lua")

addListener("post-received", function(post)
	print("Receiving post (eventstest.lua):")
	for key,value in pairs(post) do
		print(key,"=",value)
	end
end)

doEvents()