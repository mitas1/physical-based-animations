using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class ParticleSystemController : MonoBehaviour {

    ParticleSystem.EmissionModule emissionModule;

    void Start () {
        emissionModule = gameObject.GetComponent<ParticleSystem> ().emission;
    }

    public void ChangeParticlesPerSecond(Slider slider) {
        Debug.Log (slider.value);
        emissionModule.enabled = true;
        emissionModule.rateOverTime = new ParticleSystem.MinMaxCurve(slider.value);
    }
}
