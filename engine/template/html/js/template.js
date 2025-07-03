window.onload = function() {
  document.querySelectorAll('img').forEach(function(img) {
    var caption = document.createElement('figcaption');
    caption.innerHTML = "<strong>Figure: </strong>" + img.alt;
    img.parentNode.insertBefore(caption, img.nextSibling);
  });
}
