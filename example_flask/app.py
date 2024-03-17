from flask import Flask
from flask_decorator import check_resource_access

app = Flask(__name__)

app.config['APP_NAME'] = "TEST_APP"
app.config['VERIFY_URL'] = f"http://localhost:8080/api/v1/verify/{app.config['APP_NAME']}"
app.config['VERIFY_TLS'] = False

@app.route("/")
@check_resource_access(["p_user"])
def protected_index():
    return {"msg":"super secret"}


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5001, debug=True)