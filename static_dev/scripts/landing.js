$(function() {
  var sizeCallToActionExplanation;
  (sizeCallToActionExplanation = function() {
    var $items = $('.call-to-action-explanation > li');
    $items.css('height', 'auto');
    var max = 0;
    $items.each(function() {
      var thisHeight;
      if ((thisHeight = $(this).height()) > max) { max = thisHeight; }
    }).height(max);
  })();

  $(window).on('resize', sizeCallToActionExplanation);
});
