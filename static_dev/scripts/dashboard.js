window.current_user = {"id":14,"email":"matthewc267@gmail.com","last_message_id":"1471e7fae00642d4","created_at":"2014-07-31T20:31:40Z","updated_at":"2014-07-31T20:31:59Z","purchases":[{"id":103,"purchase_price":12.67,"kickback_amount":-1.35,"purchase_at":"2014-07-29T00:00:00Z","created_at":"2014-07-31T20:31:59Z","updated_at":"2014-07-31T20:32:10Z","seller_name":"Amazon.com","current_price":12.67,"was_kickbacked":false,"product":{"id":38,"productId":"B000EVPJXI","name":" BladesUSA 1806Bk Samurai Wooden Training Bokken Black 39-Inch Overall ","url":"http://www.amazon.com/dp/B000EVPJXI/ref=pe_385040_30332200_TE_item","scraped_at":"2014-07-29T00:00:00Z","created_at":"2014-07-31T20:31:59Z","updated_at":"2014-07-31T20:32:10Z","purchases":null,"small_image_url":"http://ecx.images-amazon.com/images/I/41UrmlIdj%2BL._SL75_.jpg"}},{"id":104,"purchase_price":5.44,"kickback_amount":-3.44,"purchase_at":"2014-07-29T00:00:00Z","created_at":"2014-07-31T20:31:59Z","updated_at":"2014-07-31T20:31:59Z","seller_name":"DIYMOBILITY™ (1000+ FEEDBACK/ SAME DAY USA SHIPPER. OVERNIGHT. 2 DAY AVAILABLE)","current_price":5.44,"was_kickbacked":false,"product":{"id":39,"productId":"B00CW2S3BC","name":" iPhone 5 5G Battery Replacement + TOOL SET - iPhone 5 Li-ion Battery Replacement. This Rechargeable lithium-ion polymer battery will replace your exha ","url":"http://www.amazon.com/dp/B00CW2S3BC/ref=pe_385040_30332200_TE_item","scraped_at":"2014-07-29T00:00:00Z","created_at":"2014-07-31T20:31:59Z","updated_at":"2014-07-31T20:31:59Z","purchases":null,"small_image_url":""}},{"id":105,"purchase_price":25.34,"kickback_amount":-0.46,"purchase_at":"2014-07-29T00:00:00Z","created_at":"2014-07-31T20:31:59Z","updated_at":"2014-07-31T20:32:10Z","seller_name":"Amazon.com","current_price":25.8,"was_kickbacked":false,"product":{"id":40,"productId":"B007CRUBV2","name":" Rapiddominance Classic Military Messenger Bags, Black ","url":"http://www.amazon.com/dp/B007CRUBV2/ref=pe_385040_30332200_TE_item","scraped_at":"2014-07-29T00:00:00Z","created_at":"2014-07-31T20:31:59Z","updated_at":"2014-07-31T20:32:10Z","purchases":null,"small_image_url":"http://ecx.images-amazon.com/images/I/41fpKJD5qRL._SL75_.jpg"}},{"id":102,"purchase_price":6.25,"kickback_amount":0,"purchase_at":"2014-07-29T00:00:00Z","created_at":"2014-07-31T20:31:59Z","updated_at":"2014-07-31T20:31:59Z","seller_name":"BC Novelties","current_price":6.25,"was_kickbacked":false,"product":{"id":37,"productId":"B000G2OZ30","name":" Japanese WW2 Battle Flag 3 x 5 NEW 3x5 WWII Rising Sun ","url":"http://www.amazon.com/dp/B000G2OZ30/ref=pe_385040_30332200_TE_item","scraped_at":"2014-07-29T00:00:00Z","created_at":"2014-07-31T20:31:59Z","updated_at":"2014-07-31T20:31:59Z","purchases":null,"small_image_url":""}},{"id":101,"purchase_price":6.79,"kickback_amount":0,"purchase_at":"2014-07-29T00:00:00Z","created_at":"2014-07-31T20:31:59Z","updated_at":"2014-07-31T20:31:59Z","seller_name":"BCS inc.","current_price":6.79,"was_kickbacked":false,"product":{"id":36,"productId":"B00EP95TYO","name":" 1pc Replacement End Cap Cover for Jawbone UP 2 2nd Gen 2.0 Bracelet Band Cap Dust Protector (not for the 1st Gen) ","url":"http://www.amazon.com/dp/B00EP95TYO/ref=pe_385040_30332200_TE_item","scraped_at":"2014-07-29T00:00:00Z","created_at":"2014-07-31T20:31:59Z","updated_at":"2014-07-31T20:31:59Z","purchases":null,"small_image_url":""}}]};
window.lookback = 7;
window.offset = 0;

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
  
  var ctx = $('#savings-chart').get(0).getContext('2d');
  initializePieChart(ctx);

  ctx = $('#breakdown-chart').get(0).getContext('2d');
  initializeBreakdownChart(ctx, getLookback(window.offset, window.lookback, 1));

  $('.breakdown .nav-left').on('click', function() {
    window.offset += 7;
    initializeBreakdownChart(ctx, getLookback(window.offset, window.lookback, 1));
  });

  $('.breakdown .nav-right').on('click', function() {
    window.offset -= 7;
    initializeBreakdownChart(ctx, getLookback(window.offset, window.lookback, 1));
  })
});

function initializePieChart(ctx) {
  var spent = 0;
  var saved = 0;
  for (var i = 0; i < window.current_user.purchases.length; i++) {
    spent += window.current_user.purchases[i].purchase_price;
    saved -= window.current_user.purchases[i].kickback_amount;
  }

  pieData = [
    {
        value: spent,
        color:"#FF9900",
        highlight: "#FFB31A",
    },
    {
        value: saved,
        color: "#87d1a1",
        highlight: "#A1EBBB",
    }
  ];

  $('.chart-container .total-spent .amount').text(spent);
  $('.chart-container .total-saved .amount').text(saved);

  // Get context with jQuery - using jQuery's .get() method.
  return new Chart(ctx).Doughnut(pieData,{
    responsive: true,
    percentageInnerCutout : 75,
  });

}

function initializeBreakdownChart(ctx, lookback) {
  breakdownData = {
    labels: lookback.labels,
    datasets: [
        {
            label: "Spent",
            fillColor: "rgba(255, 145, 0, 0.75)",
            strokeColor: "rgba(255, 145, 0, 1)",
            pointColor: "rgba(255, 145, 0, 1)",
            pointStrokeColor: "rgb(255, 189, 114)",
            pointHighlightFill: "rgb(255, 154, 66)",
            pointHighlightStroke: "rgba(220,220,220,1)",
            data: lookback.spent
        },
        {
            label: "Saved",
            fillColor: "rgba(138, 211, 159, 0.75)",
            strokeColor: "rgba(138, 211, 159, 1)",
            pointColor: "rgba(138, 211, 159, 1)",
            pointStrokeColor: "rgb(191, 244, 201)",
            pointHighlightFill: "rgb(181, 220, 191)",
            pointHighlightStroke: "rgba(151,187,205,1)",
            data: lookback.saved
        }
    ]
  };

  return new Chart(ctx).Line(breakdownData, {
    responsive: true,
    bezierCurve : true,
    scaleOverride: true,
    scaleSteps: 5,
    scaleStartValue: 0,
    scaleStepWidth: getScaleStepWidth(),
    scaleIntegersOnly: true,
    scaleBeginAtZero: true,
    scaleFontFamily: "'Lato', sans-serif",
    scaleFontStyle: "500",
    scaleFontSize: 30,
    scaleShowGridLines : true,
    maintainAspectRatio: false,
    showTooltips: true
  });

  function getScaleStepWidth() {
    var max = Math.max.apply(null, lookback.spent);
    if (max === 0) { max = 5}
    return Math.ceil(max / 5);
  }
}

function getLookback(daysBackStart, daysBackEnd, split) {
  var result = {
    labels: [],
    saved: [],
    spent: []
  };

  var dateRange = getDateRange(daysBackStart, daysBackEnd);
  var explodedDateRange = splitDateRange(dateRange, split);

  for (var k = 0; k < explodedDateRange.length - 1; k++) {
    result.labels.push(explodedDateRange[k].toLocaleDateString());
  }

  result.saved = Array.apply(null, new Array(result.labels.length)).map(Number.prototype.valueOf,0);
  result.spent = Array.apply(null, new Array(result.labels.length)).map(Number.prototype.valueOf,0);

  for (var i = 0; i < current_user.purchases.length; i++){
    var purchase = current_user.purchases[i];
    purchase.purchase_at = new Date(purchase.purchase_at);
    for (var j = 0; j < explodedDateRange.length - 1; j++) {
      var start = explodedDateRange[j];
      var end = explodedDateRange[j + 1];
      if (purchase.purchase_at > start && purchase.purchase_at <= end) {
        result.saved[j] -= purchase.kickback_amount;
        result.spent[j] += purchase.purchase_price;
      }
    }
  }

  return result;
}

function getDateRange(daysBackStart, lookback) {
  var endDate = new Date();
  endDate.setHours(0);
  endDate.setMinutes(0);
  endDate.setSeconds(0);
  endDate.setDate(endDate.getDate() - daysBackStart);
  if (endDate.getDay() !== 0) {
    endDate.setDate(endDate.getDate() + 7 - endDate.getDay());
  }

  var startDate = new Date(endDate);
  startDate.setDate(startDate.getDate() - lookback);

  return [startDate, endDate];
}

function splitDateRange(range, daysSplit) {
  var start = range[0];
  var end = range[1];
  var result = [start];
  while(result[result.length - 1] < end) {
    var tempDate = new Date(result[result.length - 1]);
    tempDate.setDate(tempDate.getDate() + daysSplit);
    result.push(tempDate);
  }

  return result;
}
