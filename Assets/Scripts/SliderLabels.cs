using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class SliderLabels : MonoBehaviour {

    public Text minLabel;
    public Text maxLabel;

    void Start () {
        Slider slider = gameObject.GetComponent<Slider> ();
        minLabel.text = ((int) slider.minValue).ToString();
        maxLabel.text = ((int) slider.maxValue).ToString();
    }
}
