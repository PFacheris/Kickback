'use strict';
// generated on 2014-07-18 using generator-gulp-webapp 0.1.0

var gulp = require('gulp');

// load plugins
var $ = require('gulp-load-plugins')();

gulp.task('styles', function () {
    return gulp.src('static_dev/styles/main.scss')
        .pipe($.rubySass({
            style: 'expanded',
            precision: 10,
            loadPath: ['static_dev/styles']
        }))
        .pipe(gulp.dest('.tmp/styles'))
        .pipe($.size());
});

gulp.task('scripts', function () {
    return gulp.src('static_dev/scripts/**/*.js')
        .pipe($.jshint())
        .pipe($.jshint.reporter(require('jshint-stylish')))
        .pipe($.size());
});

gulp.task('html', ['styles', 'scripts', 'wiredep'], function () {
    var jsFilter = $.filter('**/*.js');
    var cssFilter = $.filter('**/*.css');

    return gulp.src('static_dev/*/*.html')
        .pipe($.useref.assets({searchPath: '{.tmp,static_dev}'}))
        .pipe(jsFilter)
          .pipe($.sourcemaps.init())
            .pipe($.uglify())
          .pipe($.sourcemaps.write())
        .pipe(jsFilter.restore())
        .pipe(cssFilter)
          .pipe($.csso())
        .pipe(cssFilter.restore())
        .pipe($.useref.restore())
        .pipe($.useref())
        .pipe(gulp.dest('static_dist'))
        .pipe($.size());
});

gulp.task('images', function () {
    return gulp.src('static_dev/images/**/*')
        .pipe($.cache($.imagemin({
            optimizationLevel: 3,
            progressive: true,
            interlaced: true
        })))
        .pipe(gulp.dest('static_dist/images'))
        .pipe($.size());
});

gulp.task('fonts', function () {
    return $.bowerFiles({
          paths: {
            bowerDirectory: 'static_dev/bower_components',
            bowerrc: '.bowerrc',
            bowerJson: 'bower.json'
          }
        })
        .pipe($.filter('**/*.{eot,svg,ttf,woff}'))
        .pipe($.flatten())
        .pipe(gulp.dest('static_dist/fonts'))
        .pipe($.size());
});

gulp.task('extras', function () {
    return gulp.src(['static_dev/*.*', '!static_dev/*.html', '!.*'], { dot: true })
        .pipe(gulp.dest('static_dist'));
});

gulp.task('clean', function () {
    return gulp.src(['.tmp', 'static_dist'], { read: false }).pipe($.clean());
});

gulp.task('build', ['html', 'images', 'fonts', 'extras']);

gulp.task('default', ['clean'], function () {
    gulp.start('build');
});

// inject bower components
gulp.task('wiredep', function () {
    var wiredep = require('wiredep').stream;

    gulp.src('static_dev/styles/*.scss')
        .pipe(wiredep({
            directory: 'static_dev/bower_components'
        }))
        .pipe(gulp.dest('static_dev/styles'));

    gulp.src('static_dev/views/*.html')
        .pipe(wiredep({
            directory: 'static_dev/bower_components'
        }))
        .pipe(gulp.dest('.tmp/views'));
});

gulp.task('watch', function () {
    // var server = $.livereload();
    //
    // // watch for changes
    //
    // gulp.watch([
    //     'static_dev/views/*.html',
    //     '.tmp/styles/**/*.css',
    //     'static_dev/scripts/**/*.js',
    //     'static_dev/images/**/*'
    // ]).on('change', function (file) {
    //     server.changed(file.path);
    // });

    gulp.watch(['static_dev/styles/**/*.scss', 'static_dev/views/*.html', 'static_dev/scripts/**/*.js', 'bower.json'], ['html']);
    gulp.watch('static_dev/images/**/*', ['images']);
});
