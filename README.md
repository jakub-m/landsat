Simple CLI to download images from NASA Landsat images through [their API][1]


    go install -v jakub-m/landsat/landsat

    landsat -list \
        -begin 2011-01-01 -end 2018-01-01  \
        -lat 54.1734925 -lon 19.4078047  \
        -api-key $API_KEY | tee out.log

    landsat -api-key $API_KEY < out.log


[1]:https://api.nasa.gov/api.html
