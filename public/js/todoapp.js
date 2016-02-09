function ready(fn) {
  if (document.readyState != 'loading'){
    fn();
  } else {
    document.addEventListener('DOMContentLoaded', fn);
  }
}

ready(function(){
  var debounce = function(fn, delay) {
    var timeout = null
    return function() {
      var ctx = this, args = arguments;
      if (timeout) {
        clearTimeout(timeout)
      }
      timeout = setTimeout(function(){fn.apply(ctx, args);}, delay)
    }
  };

  var postForm = function (url, data, success) {
    var req = new XMLHttpRequest();
    req.open("POST", url, true);
    req.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded; charset=UTF-8');
    var pairs = []
    for (key in data) {
      pairs.push(encodeURIComponent(key)+'='+encodeURIComponent(data[key]));
    }
    req.send(pairs.join("&"));

    req.onload = function() {
      success(req);
    };
  };

  var alertError = function(error) {
    document.getElementById("error-alert").innerText = error;
  };

  var todos = document.querySelectorAll(".todo");
  for (var i = 0; i < todos.length; i++) {
    var todo = todos[i];

    var check = todo.querySelector("input[type=checkbox]");
    check.addEventListener("change",debounce(function(){
      var todo = this.parentElement;
      var check = this;
      var id = todo.getAttribute("data-id");
      postForm(
        "/task/"+id+"/update",
        { "done":(check.checked?"1":"0") },
        function(req) {
          if (req.status >= 200 && req.status < 400) {
            var data = JSON.parse(req.responseText);
            if (data["status"] == "error") {
              alertError(data["error"]);
            }
          } else {
            alertError("Ошибка сохранения");
          }
        }
      );
    }, 500));

    var label = todo.querySelector("input[type=text]");
    label.addEventListener("input",debounce(function(){
      var todo = this.parentElement;
      var label = this;
      var id = todo.getAttribute("data-id");
      postForm(
        "/task/"+id+"/update",
        { "label":label.value},
        function(req) {
          if (req.status >= 200 && req.status < 400) {
            var data = JSON.parse(req.responseText);
            if (data["status"] == "error") {
              alertError(data["error"]);
            }
          } else {
            alertError("Ошибка сохранения");
          }
        }
      );
    }, 500));

    var delbtn = todo.querySelector(".delete");
    delbtn.addEventListener("click",function(){
      var todo = this.parentElement;
      var id = todo.getAttribute("data-id");
      postForm("/task/"+id+"/destroy", null, null);
      todo.remove();
    });
  }

  var lists = document.querySelectorAll(".todolist");
  for (var i = 0; i < lists.length; i++) {
    var list = lists[i];

    var delbtn = list.querySelector(".delete");
    delbtn.addEventListener("click",function(){
      var list = this.parentElement;
      var id = list.getAttribute("data-id");
      postForm("/list/"+id+"/destroy", null, null);
      list.remove();
    });
  }

  var titles = document.querySelectorAll(".todolist-title");
  for (var i = 0; i < titles.length; i++) {
    var title = titles[i];
    title.addEventListener("input",debounce(function(){
      var title = this;
      var id = title.getAttribute("data-id");
      postForm(
        "/list/"+id+"/update",
        { "title":title.value},
        function(req) {
          if (req.status >= 200 && req.status < 400) {
            var data = JSON.parse(req.responseText);
            if (data["status"] == "error") {
              alertError(data["error"]);
            }
          } else {
            alertError("Ошибка сохранения");
          }
        }
      );
    }, 500));
  }
})
