$(function() {
  var sizeInnerChartFont;
  (sizeInnerChartFont = function() {
    var $container = $('.inner-chart-container');
    $container.find('.total-saved .amount').css('font-size', $container.height() * 0.1);
    $container.find('.total-saved .currency').css('font-size', $container.height() * 0.08);
    $container.find('.total-spent .amount').css('font-size', $container.height() * 0.07);
    $container.find('.total-spent .currency').css('font-size', $container.height() * 0.05);
    $container.find('.label').css('font-size', $container.height() * 0.04);
  })();

  $(window).on('resize', sizeInnerChartFont);

  var ctx = $("#savings-chart").get(0).getContext("2d");
  load(ctx);
});

function load(ctx) {
data = [
  {
      value: 300,
      color:"#FF9900",
      highlight: "#FFB31A",
  },
  {
      value: 50,
      color: "#87d1a1",
      highlight: "#A1EBBB",
  }
];

// Get context with jQuery - using jQuery's .get() method.
var savingsChart = new Chart(ctx).Doughnut(data,{
  responsive: true,
  percentageInnerCutout : 75,
});
}
