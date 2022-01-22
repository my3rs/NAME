window.addEventListener('load', function() {
    if (document.querySelector('#app')) {return false};

    // 计算图片宽高来获得flex的比例
    var photos = document.querySelectorAll('.photos');
    for (var i = 0; i < photos.length; i++) {
        var photos_img = photos[i].querySelectorAll('img');
        for (var j = 0; j < photos_img.length; j++) {
            photos_img[j].parentNode.parentNode.style.flex = ( photos_img[j].naturalWidth * 50 ) / photos_img[j].naturalHeight;
        }
    }

    // 给文章内没有 a 标签的 img 图片加上 a标签
    var postContentImgAll = document.querySelectorAll('.post-content p img');
    for (var j = 0; j < postContentImgAll.length; j++) {
        var postContentImg = postContentImgAll[j];
        var a = document.createElement("a");
        a.className = "p-img";
        a.href = postContentImg.src;
        var parent = postContentImg.parentNode;
        parent.appendChild(a);
        a.appendChild(postContentImg);
    }

    // 给文章内的图片a标签加data属性以及获取alt信息
    var postImgAll = document.querySelectorAll('.post-content img');
    for (var i = 0; i < postImgAll.length; i++) {
        var postImgA = postImgAll[i].parentNode;
        var postImgValue = postImgAll[i].alt;
        if (postImgValue) {
            postImgA.setAttribute('data-caption', postImgValue);
            var figcaption = document.createElement("figcaption");
            figcaption.innerHTML= `${postImgValue}`;
            postImgA.insertBefore(figcaption,postImgA.childNodes[0]);
        }
    }
});

(function postArticle() {
    var tocTitleDiv = document.createElement("div");
    tocTitleDiv.className = "title--right sulli";
    tocTitleDiv.innerText = "目录"
    // 文章目录

    var postArticles = document.querySelector(".post-articles");
    if (!postArticles) { return false; }
    var postContentH2 = document.querySelectorAll(".post-content h2");
    if (postContentH2.length) {
        postArticles.appendChild(tocTitleDiv);
    }
    var tocDiv = document.createElement("div");
    tocDiv.className = "toc";
    postArticles.appendChild(tocDiv);
    var ul = document.createElement("ul");
    tocDiv.appendChild(ul);
    for (var i = 0; i < postContentH2.length; i++) {
        var li = document.createElement("li");
        ul.appendChild(li);
        var a = document.createElement("a");
        a.href = "#cl-" + (i + 1);
        a.innerHTML = postContentH2[i].innerHTML;
        li.appendChild(a);
        postContentH2[i].setAttribute("id","cl-" + (i + 1));
    }

    window.addEventListener("scroll", function () {
        // 文章目录到评论框就消失
        var heightTop = document.scrollingElement.scrollTop;
        var commentsTop = document.querySelector('.doc_comments').offsetTop;
        heightTop > 165 ? postArticles.classList.add("is-fix") : postArticles.classList.remove("is-fix");
        heightTop < commentsTop ? postArticles.style.display = '' : postArticles.style.display = 'none';

        // 文章内部的最新文章小部件显示隐藏
        if (!document.querySelector(".doc_comments")) {return false};
        document.querySelector(".doc_comments").offsetTop < heightTop + 0
            ? document.querySelector(".js-content").classList.add("is-active")
            : document.querySelector(".js-content").classList.remove("is-active");
    });
})();

// photos 短代码实现
var photos_reg = /\[photos](.*?)\[\/photos\]/g;
var postContent = document.querySelectorAll('.post-content');
for (var o = 0; o < postContent.length; o++) {
    if (postContent[o].innerHTML.match(photos_reg)) {
        var photos_value = postContent[o].innerHTML.match(photos_reg);
        for (var i = 0; i < photos_value.length; i++) {
            photos_value[i] = photos_value[i].split(' ');
            photos_value[i].shift();
            photos_value[i].pop();
            var div = document.createElement('div');
            div.classList = 'photos';
            for (var j = 0; j < photos_value[i].length; j++) {
                photos_value[i][j] = photos_value[i][j].split(",");
                var figure = document.createElement('figure');
                figure.classList = 'photos-img';
                figure.innerHTML = `
                    <a href="${photos_value[i][j][0]}" target="_blank">
                        <img src="${photos_value[i][j][0]}!view" alt="${photos_value[i][j][1]}">
                    </a>
                    `;
                div.appendChild(figure);
            }
            postContent[o].innerHTML = postContent[o].innerHTML.replace(postContent[o].innerHTML.match(photos_reg)[0],div.outerHTML);
        }
    }
}

// 顶部固定导航栏和返回顶部到点就显示或者消失
window.addEventListener("scroll", function () {
    var heightTop = document.scrollingElement.scrollTop;
    if (heightTop > 130) {
        document.querySelector("#header-fix").style.display = "block";
        document.querySelector("#header-fix").style.position = "fixed";
        document.querySelector(".back2top").style.display = "block";
    } else {
        document.querySelector("#header-fix").style.display = "";
        document.querySelector("#header-fix").style.position = "";
        document.querySelector(".back2top").style.display = "none";
    }
});

// 夜间模式添加到本地cookie储存里
document.querySelector(".socolor").addEventListener("click", function() {
    var body = document.querySelector("body");
    if (body.getAttribute("theme") == "dark") {
        localStorage.yone_color_style = 'light';
        body.removeAttribute("theme","dark");
        document.querySelector("body").querySelector(".dark").style.display = "none";
    } else {
        localStorage.yone_color_style = 'dark';
        body.setAttribute("theme","dark");
        document.querySelector("body").querySelector(".dark").style.display = "block";
    }
});

// 滚动到顶部缓动实现
// rate表示缓动速率，默认是2
var backToTop = function (rate) {
    var doc = document.body.scrollTop? document.body : document.documentElement;
    var scrollTop = doc.scrollTop;

    var top = function () {
        scrollTop = scrollTop + (0 - scrollTop) / (rate || 8);

        // 临界判断，终止动画
        if (scrollTop <= 1) {
            doc.scrollTop = 0;
            return;
        }
        doc.scrollTop = scrollTop;
        // 动画gogogo!
        requestAnimationFrame(top);
    };
    top();
};