import json, requests
from flask import request, current_app
from functools import wraps



def check_resource_access(group=None):
    def wrapper(view_func):
        @wraps(view_func)
        def decorated(*args, **kwargs):
            token = None

            if 'Authorization' in request.headers and request.headers['Authorization'].startswith('Bearer '):
                token = request.headers['Authorization'].split(None, 1)[1].strip()

                # Check token
                r = requests.get(current_app.config['VERIFY_URL'], headers={"Authorization": f"Bearer {token}"}, verify=current_app.config['VERIFY_TLS'])

                if r.status_code == 401:
                    response_body = {'error': 'token invalid'}
                    response_body = json.dumps(response_body)
                    return response_body, 401, {'WWW-Authenticate': 'Bearer'}

                if r.status_code == 200:

                    data = r.json()

                    if group is None:
                        # Successful
                        return view_func(*args, **kwargs)
                    else:
                        # Check if it is a list
                        if isinstance(group, list):
                            for g in group:
                                if g in data['groups']:
                                    # Successful
                                    return view_func(*args, **kwargs)

                            # not successful
                            response_body = {'error': 'no access to page'}
                            response_body = json.dumps(response_body)
                            return response_body, 401, {'WWW-Authenticate': 'Bearer'}
                        else:
                            if group in data['groups']:
                                # Successful
                                return view_func(*args, **kwargs)
                            else:
                                # not successful
                                response_body = {'error': 'no access to page'}
                                response_body = json.dumps(response_body)
                                return response_body, 401, {'WWW-Authenticate': 'Bearer'}

            else:
                response_body = {'error': 'no authentication token present'}
                response_body = json.dumps(response_body)
                return response_body, 401, {'WWW-Authenticate': 'Bearer'}

        return decorated
    return wrapper





