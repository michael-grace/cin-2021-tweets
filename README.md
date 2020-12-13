# Tweet Board

###### Michael Grace, 2020

## Customising

* Needs writing

## Starting

* Copy `config.py.example` to `config.py`, and fill in the fields for your server
* Create a venv:
    * Linux:
        * `python3 -m venv venv`
        * `source venv/bin/activate`
        * `pip install -r requirements.txt`
* Start the server: `python3 tweets.py`

## Accessing

* The controller can be accessed at `/control`. Only the latest controller to connect will be usable.
* The tweet wall can be accessed at `/board`. Multiple clients can connect to this.

## Screenshots
![Tweet Board Controller](assets/control.png)