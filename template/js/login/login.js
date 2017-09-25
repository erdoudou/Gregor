/*
* @Author: Administrator
* @Date:   2017-09-22 14:35:37
* @Last Modified by:   Administrator
* @Last Modified time: 2017-09-25 13:05:41
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
bloodImg.src="/images/blood 018.png"


bloodImg.onload=function(){
	var canvas2=document.querySelector("#canvas2");
	var ctx2=canvas2.getContext("2d");
		ctx2.drawImage(bloodImg,0,0,canvas2.width,canvas2.height);

	
	// var imgData=ctx.getImageData(0,0,canvas2.width,canvas2.height)

	// console.log(imgData)
		var timer2=null;
	$("#canvas").click(function(e){
		clearInterval(timer2)
		// ctx.drawImage(bgImg,0,0,canvas.width,canvas.height)

		

		var _left=e.clientX;
		var _top=e.clientY;
//canvas2 隐藏  不能获取宽度
		var maxLeft=document.body.clientWidth-canvas2.width;
		var maxTop=document.body.clientHeight-canvas2.width;

		console.log(maxLeft,document.body.clientWidth,canvas2.width)
		if(_left>maxLeft){
			_left=maxLeft;
		}
		if(_top>maxTop){
			_top=maxTop
		}
		// console.log(_left)
		// ctx.drawImage(bloodImg,_left,_top,100,100)
		$("#canvas1").css({
			left:_left,
			top:_top
		})
		var canvas1=document.querySelector("#canvas1")
		var ctx1=canvas1.getContext("2d");
		ctx1.clearRect(0,0,canvas1.width,canvas1.height);

		
		var n=1;
				timer2=setInterval(function(){
				var imgData=ctx2.getImageData(0,0,100,n)
				ctx1.putImageData(imgData,0,0,0,0,100,100)
				// console.log(imgData,n)
				n++;
				if(n==100){
					clearInterval(timer2)
				}

			},30)
	
	})
}



// - -     -
		// var str="sadfa"
		var str="浓雾弥漫着整个森林，古老的歌谣从树林深处传出，好似呼唤着在路上的人们，陆陆续续的人跋山涉水而来，神色匆匆，消失在一片迷雾之中，渡鸦是不是的飞起，吱吱呀呀的吵着“只进不出，只进不出。。。。“，你是不是赶路的人？赶路的人是不是你？"
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
											
											$("#show span").eq(i).animate({opacity:1})
											i++
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