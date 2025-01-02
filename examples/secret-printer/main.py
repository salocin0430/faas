import os
import sys
import json
import requests

def main():
    try:
        # Get input from arguments
        input_data = sys.argv[1] if len(sys.argv) > 1 else ""
        input_json = json.loads(input_data)
        
        # Get API URL from direct inputs
        api_url = input_json.get("direct_inputs", {}).get("api_url", "")
        if not api_url:
            raise ValueError("api_url is required in direct_inputs")

        # Print all environment variables that are secrets
        print("=== Secrets ===")
        for env_var in os.environ:
            if env_var not in ["API_BASE_URL", "PATH", "PYTHONPATH"]:
                print(f"{env_var}: {os.getenv(env_var)}")

        # Make HTTP request
        print("\n=== API Response ===")
        response = requests.get(api_url)
        if response.status_code == 200:
            print(json.dumps(response.json(), indent=2))
        else:
            print(f"Error: {response.status_code}")
            
        return json.dumps({
            "secrets_found": len(os.environ) - 3,  # Exclude standard env vars
            "api_status": response.status_code
        })

    except Exception as e:
        return json.dumps({"error": str(e)})

if __name__ == "__main__":
    print(main()) 