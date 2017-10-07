using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Obstacle : MonoBehaviour {

    private Vector3 offset;
    private Vector3 original;
    private Vector3 screenPoint;
    private bool CanMove = false;

    void Start () {
        original = gameObject.transform.position;
    }

    void Update() {
        if (CanMove) {
            Vector3 curScreenPoint = new Vector3(
                Input.mousePosition.x, Input.mousePosition.y, screenPoint.z
            );
            Vector3 curPosition = Camera.main.ScreenToWorldPoint (curScreenPoint) + offset;
            gameObject.transform.position = curPosition;
        }
    }
    
    void FixedUpdate () {
        if (Input.GetMouseButtonDown (0)) {
            RaycastHit hit = new RaycastHit ();
            Ray ray = Camera.main.ScreenPointToRay (Input.mousePosition);

            if (Physics.Raycast (ray, out hit)) {
                if (hit.collider.transform.name == gameObject.transform.name) {
                    screenPoint = Camera.main.WorldToScreenPoint (gameObject.transform.position);
    
                    offset = gameObject.transform.position - Camera.main.ScreenToWorldPoint (
                        new Vector3(Input.mousePosition.x, Input.mousePosition.y, screenPoint.z)
                    );
                    CanMove = true;
                }
            }
        }

        if (Input.GetMouseButtonUp (0)) {
            CanMove = false;
        }
    }

    public void ResetPosition() {
        gameObject.transform.position = original;
    }
}
