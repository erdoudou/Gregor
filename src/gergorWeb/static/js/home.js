var canvas=document.querySelector("#canvas");
var ctx=canvas.getContext("2d");

canvas.width=document.querySelector("body").clientWidth=document.documentElement.clientWidth;
canvas.height=document.querySelector("body").clientHeight=document.documentElement.clientHeight;
console.log()
var img=new Image();
img.src="/static/img/home_bag.jpg"

img.onload=function(){
	ctx.drawImage(img,0,0,canvas.width,canvas.height)
	setTimeout(function(){
		$("#filter").animate({
			opacity:0.4
		},4000,function(){
			setInterval(function(){
				$("#filter").fadeToggle(4000,"linear")
			},1000)
		})
	},1000)
	

	//#btn
	var _li=$("#btn ul li");
	var n=0;
	// for(var )
	setInterval(function(){
		$(_li[n]).css("visibility","visible").siblings().css("visibility","none");
		n++;
		if(n==_li.length){
			n=0;
		}
	},500)
		
	// tone
	var canvas2=document.querySelector("#tone");
	var ctx2=canvas2.getContext("2d");
	var blood3=new Image();
	blood3.src="/static/img/blood3.png"

	var blood2=new Image();
	blood2.src="/static/img/blood2.png"

	var blood4=new Image();
	blood4.src="/static/img/blood4.png"

	blood3.onload=function(){
		var j=0;
	var timer2=setInterval(function(){
			ctx2.drawImage(blood3,0,j*1,100,100)
			j++;
			if(j>150){
				clearInterval(timer2)
			}
		},100)
	}
	blood2.onload=function(){
		var k=0;
		var timer3=setInterval(function(){
			ctx2.drawImage(blood2,100,k*5,100,100)
			k++;
			if(k>150){
				clearInterval(timer3)
			}
		},100)
					
	}
	blood4.onload=function(){
					ctx2.drawImage(blood4,200,0,100,100)
	}
}