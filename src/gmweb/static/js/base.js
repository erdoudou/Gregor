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
				alert(msg.message);
				window.location.href='/login';                      
            }
		}
	});
}

function lookdot(){
	
   	var playername = document.getElementById("playerid").value;	

   	if(playername == ""  ){
     	alert("用户名不能为空");
     	return false;
   	}

	$.ajax({
		async:true,
		type: "POST",
        url: "/lookdot",
        data: {playername:playername},
		dataType: "json",
		success: function(msg){
			if (msg == 1) {
				//window.open("home","_self");
            } else {
				alert("登录失败"); 
				//window.open("login","_self");                       
            }
		}
	});
}

function lookitem(){
	
   	var playername = document.getElementById("playerid").value;
	var itemtempid = document.getElementById("itemtempid").value;	

   	if(playername == ""  ){
     	alert("用户名不能为空");
     	return false;
   	}
	$("#project tbody").html("");

	$.ajax({
		async:true,
		url:'/lookitem',
		type:"POST",       
        data: {playername:playername,itemtempid:itemtempid},
		dataType:"json",
		success: function(msg){
			if (msg.code==0){
				alert("该玩家没有物品数据");
				return;
			}
			var tbody = "";
			var butDelitem='<td><a href="#" onclick="delitemuid(this);">删除</a></td>';					
			$.each(msg,function(index,value){
        		console.log("下标:"+index+"值:"+value.Itemuid);
				var trs = "";
				trs += "<tr><td>" + value.Itemuid + "</td> <td>" + value.Itemtempid + "</td>"+butDelitem+"</tr>";
				tbody += trs;
    		});
			$("#project").append(tbody);
		}
	});
}

function delitemuid(obj){
	var itemuid = $(obj).parents("tr").children("td").eq(0).text();	
	var playername = document.getElementById("playerid").value;	
  	$.ajax({
		async:true,
		url:'/lookitem',
		type:"POST",       
        data: {itemuid:itemuid,playername:playername},
		dataType:"json",
		success: function(msg){
			if (msg.code==1){
				$(obj).parents("tr").remove();
				alert("删除物品成功");
			}else{
				alert(msg.message);
			}
		}
	});
	
	
}

function addItem(){	
   	var playername = document.getElementById("playerid").value;
	var itemtempid = document.getElementById("itemtempid").value;	
  	$.ajax({
		async:true,
		url:'/additem',
		type:"POST",       
        data: {playername:playername,itemtempid:itemtempid},
		dataType:"json",
		success: function(msg){
			if (msg.code==1){
				alert("添加物品成功");
				//window.location.reload();
			}else{
				alert(msg.message);
			}
		}
	});
}

function registerUser(){
	
   	var username = document.getElementById("userid").value;
   	var password = document.getElementById("passid").value;
	var surePwd = document.getElementById("surePwd").value;
	

   	if(username == ""  ){
     	alert("用户名不能为空");
     	return false;
   	}
   	if(password == ""  ){
    	alert("密码不能为空");
     	return false;
   	}
	
	if(password != surePwd  ){
    	alert("密码不一致");
		window.location.reload(); 
     	return false;
   	} 
	
	//对页面取得的用户输入密码进行加密
	var passwordMd5 = hex_md5(password);

	$.ajax({
		async:true,
		url:'/register',
		type:"POST",       
        data: {username:username, password:passwordMd5},
		dataType:"json",
		success: function(msg){
			if (msg.code == 1) {
				alert("添加用户成功");
				window.location.href='/home';
            } else {
				alert(msg.message); 
				window.location.reload();                      
            }
		}
	});
}

function lookbox(){	
   	var playername = document.getElementById("playerid").value;	

   	if(playername == ""  ){
     	alert("用户名不能为空");
     	return false;
   	}
	$("#project tbody").html("");

	$.ajax({
		async:true,
		url:'/lookbox',
		type:"POST",       
        data: {playername:playername},
		dataType:"json",
		success: function(msg){
			if (msg.code==0){
				alert("该玩家没有物品数据");
				return;
			}
			var tbody = "";
			//var butDelitem='<td><a href="#" onclick="delboxuid(this);">删除</a></td>';					
			$.each(msg,function(index,value){
        		console.log("下标:"+index+"值:"+value.Boxuid);
				var trs = "";
				trs += "<tr><td>" + value.Boxuid + "</td> <td>" + value.Boxtempid + "</td>"+"<td>不能私自删除宝箱</td>"+"</tr>";
				tbody += trs;
    		});
			$("#project").append(tbody);
		}
	});
}

function delboxuid(obj){
	var boxuid = $(obj).parents("tr").children("td").eq(0).text();	
	var playername = document.getElementById("playerid").value;	
  	$.ajax({
		async:true,
		url:'/lookbox',
		type:"POST",       
        data: {boxuid:boxuid,playername:playername},
		dataType:"json",
		success: function(msg){
			if (msg.code==1){
				$(obj).parents("tr").remove();
				alert("删除物品成功");
			}else{
				alert(msg.message);
			}
		}
	});	
}

function lookdot(){	
   	var playername = document.getElementById("playerid").value;	

   	if(playername == ""  ){
     	alert("用户名不能为空");
     	return false;
   	}

	$.ajax({
		async:true,
		url:'/lookdot',
		type:"POST",       
        data: {playername:playername},
		dataType:"json",
		success: function(msg){
			if (msg.code==0){
				alert("该玩家没有物品数据");
				return;
			}
			console.log(msg.dot);
			$("#playerdot").val(msg.dot);
		}
	});
}


function changedot(){	
   	var playername = document.getElementById("playerid").value;
	var playerdot = document.getElementById("playerdot").value;	

   	if(playername == ""||playerdot == ""){
     	alert("用户名不能为空");
     	return false;
   	}

	$.ajax({
		async:true,
		url:'/lookdot',
		type:"POST",       
        data: {playername:playername,playerdot:playerdot},
		dataType:"json",
		success: function(msg){
			if (msg.code==0){
				alert(msg.message);
				return;
			}
			alert("修改玩家科技点成功");
		}
	});
}
