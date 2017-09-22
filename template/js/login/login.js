/*
* @Author: Administrator
* @Date:   2017-09-22 14:35:37
* @Last Modified by:   Administrator
* @Last Modified time: 2017-09-22 16:03:39
*/
$(function(){

			// 图片
var bgImg=new Image();
bgImg.src="/images/003_login.png";
var canvas=document.querySelector("canvas")
var ctx=canvas.getContext("2d");
canvas.width=document.body.clientWidth;
canvas.height=document.body.clientHeight;

bgImg.onload=function(){
	ctx.drawImage(bgImg,0,0,canvas.width,canvas.height);
}

var bloodImg=new Image();
bloodImg.src="/images/blood.png"


bloodImg.onload=function(){
	console.log("fdsaf")
	$("#canvas").click(function(e){
		ctx.clearRect(0,0,canvas.width,canvas.height);
		ctx.drawImage(bgImg,0,0,canvas.width,canvas.height)
		console.log("dsfas")
		var _left=e.clientX;
		var _top=e.clientY;
		console.log(_left)
		ctx.drawImage(bloodImg,_left,_top,100,100)
	})
}



// - -     -
		// var str="sadfa"
		var str="hxj是傻逼a sdf dsfa sdfasdfa sdfasdf assadaf ffffffffffffffffff fffffffff fsdf adsf asd fasdf asdfasdfds afdas  ffffff ffffff ffffffff"
		for(var i=0;i<str.length;i++){
			var span=$("<span>"+str.charAt(i)+"</span>");
			span.appendTo($("#show"));
		}
		var n=0;
		var flag=0;
		var timer=setInterval(function(){
			$("span").eq(n).animate({
				opacity:1
			},200)
			n++;
			if(n==$("span").length){
				clearInterval(timer)
				// console.log("hello")
				setTimeout(function(){
					$("#yes").animate({
						opacity:1
					},function(){
						$("#yes").click(function(){
							// 清除画布
							ctx.clearRect(0,0,canvas.width,canvas.height)
							ctx.drawImage(bgImg,0,0,canvas.width,canvas.height);
							

							$(".text").css("display","none")
							$("#login-box").css("display","block")
							$("#login-box").animate({opacity:1},1000,function(){
								$("#close").click(function(){
									$("#login-box").animate({
										opacity:0
									},function(){
										$("#login-box").css("display","none")
										$(".text").css("display","block")

										$("#show span").css({opacity:0})
										$(".text button").css({opacity:0});

										var i=0;
							var timer1=setInterval(function(){
											i++
											$("#show span").eq(i).animate({opacity:1})
											if(i==$("#show span").length){
												setTimeout(function(){
													$("#yes").animate({opacity:1},function(){
													$("#no").animate({opacity:1})
												})},50)
												
												clearInterval(timer1)
											}
										},100)
										

										
									})
								})
							});

						})
					})
				},300)
				setTimeout(function(){
					$("#no").animate({
						opacity:1
					},function(){
						$("#no").click(function(){
							alert("fweafwa")
						})
					})
				},500)
			}
			
		},100)






})