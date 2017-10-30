function checkUser(){
   	var username = document.getElementById("userid").value;
   	var password = document.getElementById("passid").value;
	

   	if(username == ""  ){
     	alert("用户名不能为空");
     	return false;
   	}
   	if(password == ""  ){
    	alert("密码不能为空");
     	return false;
   	}
	
	//对页面取得的用户输入密码进行加密
	var passwordMd5 = hex_md5(password);

	$.ajax({
		async:true,
		url:'/login',
		type:"POST",       
        data: {username:username, password:passwordMd5},
		dataType:"json",
		success: function(msg){
			if (msg.code == 1) {
				window.location.href='/home';
            } else {
				window.location.href='/login';                      
            }
		}
	});
}