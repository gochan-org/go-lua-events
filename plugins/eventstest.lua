print("Testing post-received event")

addListener("post-received", function(post)
	print("Receiving post:")
	for key,value in pairs(post) do
		print(key,"=",value)
	end
end)

doEvents()